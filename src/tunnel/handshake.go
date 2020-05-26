package tunnel

import (
    "errors"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "net"
)

const HandshakeProtocolVersion = 0x02

// HandshakeProtocol 握手协议
type HandshakeProtocol struct {
    Conn net.Conn

    Username string
    Password string

    Code byte
}

// Connection 发送连接消息
func (h *HandshakeProtocol) Connection() error {
    password := utils.EncryptPassword(h.Password)

    dataLength := 1 + 1 + len(h.Username) + 1 + len(password)

    data := make([]byte, dataLength)
    data[0] = HandshakeProtocolVersion
    data[1] = byte(len(h.Username))

    index := 2
    for _, d := range []byte(h.Username) {
        data[index] = d
        index++
    }

    data[index] = byte(len(password))
    index++

    for _, d := range []byte(password) {
        data[index] = d
        index++
    }

    n, err := h.Conn.Write(data)
    if n != dataLength || nil != err {
        return errors.New("写入数据失败")
    }

    return nil
}

// ReceiveC 接收连接消息
func (h *HandshakeProtocol) ReceiveC() error {
    ver := make([]byte, 1)
    n, err := h.Conn.Read(ver)
    if n != 1 || nil != err {
        return errors.New("读取版本号失败")
    }

    if HandshakeProtocolVersion != ver[0] {
        return errors.New("不支持的协议版本")
    }

    uLen := make([]byte, 1)
    n, err = h.Conn.Read(uLen)
    if n != 1 || nil != err {
        return errors.New("读取用户名称长度失败")
    }

    user := make([]byte, uLen[0])
    n, err = h.Conn.Read(user)
    if n != int(uLen[0]) || nil != err {
        return errors.New("读取用户名称失败")
    }
    h.Username = string(user)

    pLen := make([]byte, 1)
    n, err = h.Conn.Read(pLen)
    if n != 1 || nil != err {
        return errors.New("读取用户密码长度失败")
    }

    pass := make([]byte, pLen[0])
    n, err = h.Conn.Read(pass)
    if n != int(pLen[0]) || nil != err {
        return errors.New("读取用户密码失败")
    }
    h.Password = string(pass)

    return nil
}

// Send 发送响应
func (h *HandshakeProtocol) Send(code byte) error {
    data := make([]byte, 2)

    data[0] = HandshakeProtocolVersion
    data[1] = code

    n, err := h.Conn.Write(data)
    if n != 2 || nil != err {
        return errors.New("写入数据失败")
    }

    return nil
}

// ReceiveR 接收响应消息
func (h *HandshakeProtocol) ReceiveR() error {
    ver := make([]byte, 1)
    n, err := h.Conn.Read(ver)
    if n != 1 || nil != err {
        return errors.New("读取应答版本号失败")
    }

    if HandshakeProtocolVersion != ver[0] {
        return errors.New("不支持的协议版本")
    }

    code := make([]byte, 1)
    n, err = h.Conn.Read(code)
    if n != 1 || nil != err {
        return errors.New("读取响应码失败")
    }
    h.Code = code[0]

    return nil
}
