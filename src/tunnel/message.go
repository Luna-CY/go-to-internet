package tunnel

import (
    "encoding/binary"
    "errors"
    "net"
)

// NewConnectMessage 建立一个新的连接消息
func NewConnectMessage(conn net.Conn, ipType byte, dstIp string, dstPort int) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Cmd: CmdNewConnect, Code: MessageCodeNotSet, IpType: ipType, DstIp: dstIp, DstPort: dstPort}
}

// NewDataMessage 建立一个新的数据消息
func NewDataMessage(conn net.Conn, data []byte) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Cmd: CmdData, Code: MessageCodeNotSet, Data: data}
}

// NewOverMessage 建立一个结束消息
func NewOverMessage(conn net.Conn) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Cmd: CmdOver, Code: MessageCodeNotSet}
}

// NewEmptyMessage 建立一个空消息
func NewEmptyMessage(conn net.Conn) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Code: MessageCodeNotSet}
}

// MessageProtocol 消息协议
type MessageProtocol struct {
    Conn net.Conn
    Cmd  byte
    Code byte

    IpType  byte
    DstIp   string
    DstPort int

    Data []byte
}

// Send 发送消息
func (m *MessageProtocol) Send() error {
    sendData := m.getData()
    dataLength := 1 + 1 + 1 + 1 + len(sendData)

    data := make([]byte, dataLength)
    data[0] = MessageProtocolVersion
    data[1] = m.Cmd
    data[2] = m.Code
    data[3] = byte(len(sendData))

    index := 4
    for _, d := range sendData {
        data[index] = d
        index++
    }

    n, err := m.Conn.Write(data)
    if n != dataLength || nil != err {
        return errors.New("写入数据失败")
    }

    return nil
}

// Receive 接收消息
func (m *MessageProtocol) Receive() error {
    ver := make([]byte, 1)
    n, err := m.Conn.Read(ver)
    if n != 1 || nil != err {
        return errors.New("读取应答版本号失败")
    }

    if MessageProtocolVersion != ver[0] {
        return errors.New("不支持的协议版本")
    }

    cmd := make([]byte, 1)
    n, err = m.Conn.Read(cmd)
    if n != 1 || nil != err {
        return errors.New("读取指令失败")
    }
    m.Cmd = cmd[0]

    code := make([]byte, 1)
    n, err = m.Conn.Read(code)
    if n != 1 || nil != err {
        return errors.New("读取响应码失败")
    }
    m.Code = code[0]

    dataLen := make([]byte, 1)
    n, err = m.Conn.Read(dataLen)
    if n != 1 || nil != err {
        return errors.New("读取数据长度失败")
    }

    data := make([]byte, dataLen[0])
    n, err = m.Conn.Read(data)
    if n != int(dataLen[0]) || nil != err {
        return errors.New("读取数据失败")
    }
    m.Data = data

    return nil
}

// ParseDst 从data中解析目标信息
func (m *MessageProtocol) ParseDst() error {
    ipType := m.Data[0]

    m.DstPort = int(m.Data[1])<<8 | int(m.Data[2])

    switch ipType {
    case 0x01:
        m.DstIp = string(m.Data[3:])
    case 0x03:
        m.DstIp = string(m.Data[3:])
    case 0x04:
        m.DstIp = net.IP(m.Data[3:]).String()
    }

    return nil
}

// getData 组装发送时的消息内容
func (m *MessageProtocol) getData() []byte {
    switch m.Cmd {
    case CmdNewConnect:
        // 响应消息不需要携带数据
        if MessageCodeSuccess == m.Code {
            return make([]byte, 0)
        }

        data := make([]byte, 1+2+1+len(m.DstIp))

        data[0] = m.IpType
        index := 1

        bs := make([]byte, 2)
        binary.BigEndian.PutUint16(bs, uint16(m.DstPort))

        for _, d := range bs {
            data[index] = d
            index++
        }

        data[index] = byte(len(m.DstIp))
        index++

        for _, d := range []byte(m.DstIp) {
            data[index] = d
            index++
        }

        return data
    default:
        return make([]byte, 0)
    }
}
