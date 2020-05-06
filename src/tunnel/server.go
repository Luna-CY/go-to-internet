package tunnel

import (
    "errors"
    "fmt"
    "io"
    "net"
    "time"
)

// NewServer 新建一个隧道的服务端
func NewServer(src net.Conn) (*Server, error) {
    server := &Server{src: src}
    if !server.checkConnection() {
        return nil, errors.New("验证连接失败")
    }

    return server, nil
}

// Server 隧道的服务端结构体
type Server struct {
    src net.Conn
}

// Bind 双向绑定客户端以及目标服务器
func (s *Server) Bind() error {
    ip, port, err := s.receiveTarget()
    if nil != err {
        return err
    }

    dst, err := net.Dial("tcp", fmt.Sprintf("%v:%d", ip, port))
    if nil != err {
        return err
    }

    defer dst.Close()

    go func() {
        _, _ = io.Copy(s.src, dst)

    }()
    _, _ = io.Copy(dst, s.src)

    return nil
}

// checkConnection 检查连接是否是私有协议
func (s *Server) checkConnection() bool {
    _, _, err := s.receiveUserInfo()
    if nil != err {
        return false
    }

    // TODO: 检查用户信息是否有效

    if err = s.sendTimeout(); nil != err {
        return false
    }

    return true
}

// receiveUserInfo 获取用户信息
func (s *Server) receiveUserInfo() (string, string, error) {
    ver := make([]byte, 1)
    n, err := s.src.Read(ver)
    if n != 1 || nil != err {
        return "", "", errors.New("读取版本号失败")
    }

    uLen := make([]byte, 1)
    n, err = s.src.Read(uLen)
    if n != 1 || nil != err {
        return "", "", errors.New("读取用户名称长度失败")
    }

    user := make([]byte, uLen[0])
    n, err = s.src.Read(user)
    if n != int(uLen[0]) || nil != err {
        return "", "", errors.New("读取用户名称失败")
    }

    pLen := make([]byte, 1)
    n, err = s.src.Read(pLen)
    if n != 1 || nil != err {
        return "", "", errors.New("读取用户密码长度失败")
    }

    pass := make([]byte, pLen[0])
    n, err = s.src.Read(pass)
    if n != int(pLen[0]) || nil != err {
        return "", "", errors.New("读取用户密码失败")
    }

    return string(user), string(pass), nil
}

// sendTimeout 发送超时时间设置
func (s *Server) sendTimeout() error {
    dataLength := 11

    data := make([]byte, dataLength)
    data[0] = VER01

    index := 1
    for _, d := range []byte(fmt.Sprintf("%d", time.Now().Unix())) {
        data[index] = d
        index++
    }

    n, err := s.src.Write(data)
    if n != dataLength || nil != err {
        s.src.Close()

        return errors.New("写入数据失败")
    }

    return nil
}

// ReceiveTarget 获取目标服务器信息
func (s *Server) receiveTarget() (string, int, error) {
    ver := make([]byte, 1)
    n, err := s.src.Read(ver)
    if n != 1 || nil != err {
        return "", 0, errors.New("读取版本号失败")
    }

    port := make([]byte, 2)
    n, err = s.src.Read(port)
    if n != 2 || nil != err {
        return "", 0, errors.New("读取端口号失败")
    }

    ipType := make([]byte, 1)
    n, err = s.src.Read(ipType)
    if n != 1 || nil != err {
        return "", 0, errors.New("读取ip类型失败")
    }

    ipLen := make([]byte, 1)
    n, err = s.src.Read(ipLen)
    if n != 1 || nil != err {
        return "", 0, errors.New("读取ip长度失败")
    }

    ip := make([]byte, ipLen[0])
    n, err = s.src.Read(ip)
    if n != int(ipLen[0]) || nil != err {
        return "", 0, errors.New("读取ip失败")
    }

    err = s.sendRes()
    if nil != err {
        return "", 0, err
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

    return ipString, int(port[0])<<8 | int(port[1]), nil
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

    n, err := s.src.Write(data)
    if n != dataLength || nil != err {
        s.src.Close()
        return errors.New("写入数据失败")
    }

    return nil
}
