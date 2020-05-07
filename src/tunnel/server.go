package tunnel

import (
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "io"
    "net"
)

// NewServer 新建一个隧道的服务端
func NewServer(src net.Conn) (*Server, error) {
    server := &Server{clientConn: src}
    if !server.checkConnection() {
        return nil, errors.New("验证连接失败")
    }

    return server, nil
}

// Server 隧道的服务端结构体
type Server struct {
    clientConn net.Conn

    dstIp   string
    dstPort int
}

// Bind 双向绑定客户端以及目标服务器
func (s *Server) Bind() error {
    fmt.Printf("建立请求 -> %v:%d\n", s.dstIp, s.dstPort)

    dst, err := net.Dial("tcp", fmt.Sprintf("%v:%d", s.dstIp, s.dstPort))
    if nil != err {
        return err
    }
    defer dst.Close()

    go func() {
        defer s.clientConn.Close()
        _, _ = io.Copy(s.clientConn, dst)

    }()
    _, _ = io.Copy(dst, s.clientConn)

    return nil
}

// checkConnection 检查连接是否是私有协议
func (s *Server) checkConnection() bool {
    _, _, err := s.receiveUserInfo()
    if nil != err {
        logger.Debugf("解析协议失败: %v", err)

        return false
    }

    // TODO: 检查用户信息是否有效

    if err := s.parseTarget(); nil != err {
        logger.Errorf("解析目标数据失败: %v", err)

        return false
    }

    if err := s.sendRes(); nil != err {
        logger.Error("发送协议响应数据失败")

        return false
    }

    return true
}

// receiveUserInfo 获取用户信息
func (s *Server) receiveUserInfo() (string, string, error) {
    ver := make([]byte, 1)
    n, err := s.clientConn.Read(ver)
    if n != 1 || nil != err {
        return "", "", errors.New("读取版本号失败")
    }

    if VER01 != ver[0] {
        return "", "", errors.New("不支持的协议版本")
    }

    uLen := make([]byte, 1)
    n, err = s.clientConn.Read(uLen)
    if n != 1 || nil != err {
        return "", "", errors.New("读取用户名称长度失败")
    }

    user := make([]byte, uLen[0])
    n, err = s.clientConn.Read(user)
    if n != int(uLen[0]) || nil != err {
        return "", "", errors.New("读取用户名称失败")
    }

    pLen := make([]byte, 1)
    n, err = s.clientConn.Read(pLen)
    if n != 1 || nil != err {
        return "", "", errors.New("读取用户密码长度失败")
    }

    pass := make([]byte, pLen[0])
    n, err = s.clientConn.Read(pass)
    if n != int(pLen[0]) || nil != err {
        return "", "", errors.New("读取用户密码失败")
    }

    return string(user), string(pass), nil
}

// parseTarget 解析目标信息
func (s *Server) parseTarget() error {
    port := make([]byte, 2)
    n, err := s.clientConn.Read(port)
    if n != 2 || nil != err {
        return errors.New("解析端口失败")
    }
    s.dstPort = int(port[0])<<8 | int(port[1])

    ipType := make([]byte, 1)
    n, err = s.clientConn.Read(ipType)
    if n != 1 || nil != err {
        return errors.New("解析ip类型失败")
    }

    ipLen := make([]byte, 1)
    n, err = s.clientConn.Read(ipLen)
    if n != 1 || nil != err {
        return errors.New("解析ip长度失败")
    }

    ip := make([]byte, ipLen[0])
    n, err = s.clientConn.Read(ip)
    if n != int(ipLen[0]) || nil != err {
        return errors.New("解析ip失败")
    }

    var ipString string
    switch ipType[0] {
    case 0x01:
        ipString = net.IP(ip[0:4]).String()
    case 0x03:
        ipString = string(ip)
    case 0x04:
        ipString = net.IP(ip[0:16]).String()
    }
    s.dstIp = ipString

    return nil
}

// sendRes 发送响应数据
func (s *Server) sendRes() error {
    dataLength := 1 + 1 + 1 + len("OK")
    data := make([]byte, dataLength)
    data[0] = VER01
    data[1] = 0x00
    data[2] = byte(len("OK"))

    index := 3
    for _, d := range []byte("OK") {
        data[index] = d
        index++
    }

    n, err := s.clientConn.Write(data)
    if n != dataLength || nil != err {
        s.clientConn.Close()

        return errors.New("写入数据失败")
    }

    return nil
}
