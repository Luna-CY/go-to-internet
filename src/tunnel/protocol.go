package tunnel

import (
    "errors"
    "net"
)

// 握手协议
//
// 建立连接
// VER USER_L USER PASS_L PASS
//  1    1     N     1     N
//
// 响应消息
// VER CODE
//  1   1
//
// 通信协议
//
// VER CMD RID DATA_L DATA
//  1   1   8    1     N

const Ver = 0x02

const Success = 0x00
const SuccessMessage = "OK"

const ConnectionUpperLimit = 0x01
const ConnectionUpperLimitMessage = "已到达连接上限"

const CheckConnectTargetIp = "0.0.0.0"
const CheckConnectTargetPort = 0

// Protocol v1协议
type Protocol struct {
    Conn net.Conn

    username string
    password string
    dstIp    string
    dstPort  int
}

// receiveUserInfo 接收用于验证的用户信息
func (p *Protocol) Receive() error {
    ver := make([]byte, 1)
    n, err := p.Conn.Read(ver)
    if n != 1 || nil != err {
        return errors.New("读取版本号失败")
    }

    if Ver != ver[0] {
        return errors.New("不支持的协议版本")
    }

    uLen := make([]byte, 1)
    n, err = p.Conn.Read(uLen)
    if n != 1 || nil != err {
        return errors.New("读取用户名称长度失败")
    }

    user := make([]byte, uLen[0])
    n, err = p.Conn.Read(user)
    if n != int(uLen[0]) || nil != err {
        return errors.New("读取用户名称失败")
    }
    p.username = string(user)

    pLen := make([]byte, 1)
    n, err = p.Conn.Read(pLen)
    if n != 1 || nil != err {
        return errors.New("读取用户密码长度失败")
    }

    pass := make([]byte, pLen[0])
    n, err = p.Conn.Read(pass)
    if n != int(pLen[0]) || nil != err {
        return errors.New("读取用户密码失败")
    }
    p.password = string(pass)

    return nil
}

// receiveDstInfo 接收连接目标信息
func (p *Protocol) receiveDstInfo() error {
    port := make([]byte, 2)
    n, err := p.Conn.Read(port)
    if n != 2 || nil != err {
        return errors.New("解析端口失败")
    }
    p.dstPort = int(port[0])<<8 | int(port[1])

    ipType := make([]byte, 1)
    n, err = p.Conn.Read(ipType)
    if n != 1 || nil != err {
        return errors.New("解析ip类型失败")
    }

    ipLen := make([]byte, 1)
    n, err = p.Conn.Read(ipLen)
    if n != 1 || nil != err {
        return errors.New("解析ip长度失败")
    }

    ip := make([]byte, ipLen[0])
    n, err = p.Conn.Read(ip)
    if n != int(ipLen[0]) || nil != err {
        return errors.New("解析ip失败")
    }

    var ipString string
    switch ipType[0] {
    case 0x01:
        ipString = string(ip)
    case 0x03:
        ipString = string(ip)
    case 0x04:
        ipString = net.IP(ip[0:16]).String()
    }
    p.dstIp = ipString

    return nil
}

// Send 发送响应
func (p *Protocol) Send(code byte) error {
    data := make([]byte, 2)

    data[0] = Ver
    data[1] = code

    n, err := p.Conn.Write(data)
    if n != 2 || nil != err {
        return errors.New("写入数据失败")
    }

    return nil
}

// GetUsername username getter
func (p *Protocol) GetUsername() string {
    return p.username
}

// GetPassword password getter
func (p *Protocol) GetPassword() string {
    return p.password
}
