package tunnel

import "net"

const MessageProtocolVersion = 0x01

const CmdNewConnect = 0x01
const CmdData = 0x02

// NewConnectMessage 建立一个新的连接消息
func NewConnectMessage(conn net.Conn, dstIp string, dstPort int) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Cmd: CmdNewConnect, DstIp: dstIp, DstPort: dstPort}
}

// NewDataMessage 建立一个新的数据消息
func NewDataMessage(conn net.Conn, data []byte) *MessageProtocol {
    return &MessageProtocol{Conn: conn, Cmd: CmdData, Data: data}
}

// MessageProtocol 消息协议
type MessageProtocol struct {
    Conn net.Conn
    Cmd  byte

    DstIp   string
    DstPort int

    Data []byte
}

// Send 发送消息
func (m *MessageProtocol) Send() error {
    return nil
}

// Receive 接收消息
func (m *MessageProtocol) Receive() error {
    return nil
}
