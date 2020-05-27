package socket

import (
    "encoding/binary"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "net"
)

// Socket 结构体
type Socket struct {
    Hostname  string // 服务器域名
    Port      int    // 服务器端口
    LocalAddr string // 本地监听地址
    LocalPort int    // 本地监听端口

    Username string // 身份认证用户名
    Password string // 身份认证密码

    Verbose bool // 详细模式
}

// Start 启动本地服务监听
func (s *Socket) Start() {
    logger.Infof("启动监听 %v:%d ...\n", s.LocalAddr, s.LocalPort)

    listen, err := net.Listen("tcp", fmt.Sprintf("%v:%d", s.LocalAddr, s.LocalPort))
    if nil != err {
        logger.Errorf("监听地址 %v:%d 失败", s.LocalAddr, s.LocalPort)

        return
    }
    defer listen.Close()

    client := &client{Socket: s}
    if err := client.Init(); nil != err {
        logger.Errorf("初始化代理客户端失败: %v", err)

        return
    }

    for {
        conn, err := listen.Accept()
        if nil != err {
            continue
        }

        go client.Accept(conn)
    }
}

// isSocks5 检查连接是否是socks5协议
func (s *Socket) isSocks5(conn net.Conn) bool {
    buffer := make([]byte, 1)
    _, _ = conn.Read(buffer)

    return 0x05 == buffer[0]
}

// authorize 身份验证
func (s *Socket) authorize(conn net.Conn) bool {
    // 身份验证应答：不需要验证
    _, _ = conn.Write([]byte{0x05, 0x00})

    // 这里curl客户端好像会给一个[2 0 1]的应答数据流，先忽略掉
    buffer := make([]byte, 3)
    _, _ = conn.Read(buffer)

    return true
}

// isConnectCmd 检查连接是否是connect类型的命令
func (s *Socket) isConnectCmd(conn net.Conn) bool {
    // 取出前三个字节
    buffer := make([]byte, 3)
    _, _ = conn.Read(buffer)

    if 0x05 != buffer[0] {
        return false
    }

    return 0x01 == buffer[1]
}

// getRemoteAddr 获取远程目标的ip和端口
func (s *Socket) getRemoteAddr(conn net.Conn) (byte, string, int, error) {
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

        buffer := make([]byte, int(length[0])+2)
        n, _ := conn.Read(buffer)

        ip = string(buffer[:n-2])
        port = int(buffer[n-2])<<8 | int(buffer[n-1])
    case 0x04:
        buffer := make([]byte, 18)
        n, _ := conn.Read(buffer)

        ip = net.IP(buffer[0:16]).String()
        port = int(buffer[n-2])<<8 | int(buffer[n-1])
    default:
        return 0x00, "", 0, errors.New("解析数据失败")
    }

    return atyp[0], ip, port, nil
}

// sendAck 发送socket确认信息
func (s *Socket) sendAck(src net.Conn, response byte) {
    ack := make([]byte, 4+1+len(s.Hostname)+2)
    ack[0] = 0x05     // VER
    ack[1] = response // REP
    ack[2] = 0x00     // RSV
    ack[3] = 0x03     // ATYP: 域名
    ack[4] = byte(len(s.Hostname))

    index := 5
    for _, d := range []byte(s.Hostname) {
        ack[index] = d
        index++
    }

    bs := make([]byte, 2)
    binary.BigEndian.PutUint16(bs, uint16(s.Port))
    for _, d := range bs {
        ack[index] = d
        index++
    }

    _, _ = src.Write(ack)
}
