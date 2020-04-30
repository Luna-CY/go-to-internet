package socket

import (
    "encoding/binary"
    "errors"
    "fmt"
    "log"
    "net"
)

// Socket 结构体
type Socket struct {
    Hostname  string // 服务器域名
    Port      int    // 服务器端口
    LocalAddr string // 本地监听地址
    LocalPort int    // 本地监听端口

    Authorize bool   // 是否需要身份认证
    Username  string // 身份认证用户名
    Password  string // 身份认证密码
}

// Start 启动本地服务监听
func (c *Socket) Start() {
    fmt.Printf("启动监听 %v:%d ...\n", c.LocalAddr, c.LocalPort)

    listen, err := net.Listen("tcp", fmt.Sprintf("%v:%d", c.LocalAddr, c.LocalPort))
    if nil != err {
        log.Fatal(fmt.Sprintf("监听地址 %v:%d 失败", c.LocalAddr, c.LocalPort))
    }
    defer listen.Close()

    for {
        conn, err := listen.Accept()
        if nil != err {
            log.Fatal("接收请求失败")
        }

        go c.connection(conn)
    }
}

// connection 处理socket连接请求
func (c *Socket) connection(conn net.Conn) {
    defer conn.Close()

    if !c.isSocks5(conn) || !c.authorize(conn) {
        return
    }

    if !c.isConnectCmd(conn) {
        _, _ = conn.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

        return
    }

    ip, port, err := c.getRemoteAddr(conn)
    if nil != err {
        _, _ = conn.Write([]byte{0x05, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

        return
    }

    ack := make([]byte, 4+1+len(c.Hostname)+2)
    ack[0] = 0x05 // VER
    ack[1] = 0x00 // REP
    ack[2] = 0x00 // RSV
    ack[3] = 0x03 // ATYP: 域名
    ack[4] = byte(len(c.Hostname))

    index := 5
    for _, d := range []byte(c.Hostname) {
        ack[index] = d
        index++
    }

    bs := make([]byte, 2)
    binary.BigEndian.PutUint16(bs, uint16(c.Port))
    for _, d := range bs {
        ack[index] = d
        index++
    }

    _, _ = conn.Write(ack)

    // 处理http请求
    http := HTTP{Sock: c, Conn: &conn, TargetIp: ip, TargetPort: port}
    http.request()
}

// isSocks5 检查连接是否是socks5协议
func (c *Socket) isSocks5(conn net.Conn) bool {
    buffer := make([]byte, 1)
    _, _ = conn.Read(buffer)

    return 0x05 == buffer[0]
}

// authorize 身份验证
func (c *Socket) authorize(conn net.Conn) bool {
    // 身份验证应答：不需要验证
    _, _ = conn.Write([]byte{0x05, 0x00})

    // 这里curl客户端好像会给一个[2 0 1]的应答数据流，先忽略掉
    buffer := make([]byte, 3)
    _, _ = conn.Read(buffer)

    return true
}

// isConnectCmd 检查连接是否是connect类型的命令
func (c *Socket) isConnectCmd(conn net.Conn) bool {
    // 取出前三个字节
    buffer := make([]byte, 3)
    _, _ = conn.Read(buffer)

    if 0x05 != buffer[0] {
        return false
    }

    return 0x01 == buffer[1]
}

// getRemoteAddr 获取远程目标的ip和端口
func (c *Socket) getRemoteAddr(conn net.Conn) (string, int, error) {
    atyp := make([]byte, 1)
    _, _ = conn.Read(atyp)

    var ip, port = "", 0

    switch atyp[0] {
    case 0x01:
        buffer := make([]byte, 6)
        n, _ := conn.Read(buffer)

        ip = net.IP(buffer[0:4]).String()
        port = int(buffer[n-2])<<8 | int(buffer[n-1])
    case 0x03:
        length := make([]byte, 1)
        _, _ = conn.Read(length)

        buffer := make([]byte, int(length[0]))
        n, _ := conn.Read(buffer)

        ip = string(buffer[2 : n-2])
    case 0x04:
        buffer := make([]byte, 18)
        n, _ := conn.Read(buffer)

        ip = net.IP(buffer[0:16]).String()
        port = int(buffer[n-2])<<8 | int(buffer[n-1])
    default:
        return "", 0, errors.New("解析数据失败")
    }

    return ip, port, nil
}
