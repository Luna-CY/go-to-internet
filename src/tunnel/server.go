package tunnel

import (
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/time/rate"
    "io"
    "net"
)

// NewServer 新建一个隧道的服务端
func NewServer(src net.Conn, userConfig *config.UserConfig, verbose bool) (*Server, error) {
    server := &Server{clientConn: src, userConfig: userConfig, verbose: verbose}
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

    userConfig *config.UserConfig
    userInfo   *config.UserInfo

    verbose bool
}

// Bind 双向绑定客户端以及目标服务器
func (s *Server) Bind() error {
    // 如果是验证连接的请求直接返回完成
    if CheckConnectTargetIp == s.dstIp && CheckConnectTargetPort == s.dstPort {
        return nil
    }

    logger.Infof("建立连接请求 -> %v:%d\n", s.dstIp, s.dstPort)

    dst, err := net.Dial("tcp", fmt.Sprintf("%v:%d", s.dstIp, s.dstPort))
    if nil != err {
        return err
    }
    defer dst.Close()

    var limiter *rate.Limiter
    if 0 != s.userInfo.MaxRate {
        limiter = rate.NewLimiter(rate.Limit(s.userInfo.MaxRate*1024), s.userInfo.MaxRate*512/2)
    }

    go func() {
        defer s.clientConn.Close()
        if _, err := utils.Copy(s.clientConn, dst, limiter); nil != err && io.EOF != err && s.verbose {
            logger.Errorf("目标服务器 -> 隧道客户端: %v", err)
        }
    }()

    if _, err := utils.Copy(dst, s.clientConn, limiter); nil != err && io.EOF != err && s.verbose {
        logger.Errorf("隧道客户端 -> 目标服务器: %v", err)
    }

    return nil
}

// checkConnection 检查连接是否是私有协议
func (s *Server) checkConnection() bool {
    user, pass, err := s.receiveUserInfo()
    if nil != err {
        if s.verbose {
            logger.Debugf("解析协议失败: %v", err)
        }

        return false
    }

    userInfo, ok := s.userConfig.Users[user]
    if !ok {
        return false
    }

    err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(pass))
    if nil != err {
        return false
    }
    s.userInfo = userInfo

    if err := s.parseTarget(); nil != err {
        if s.verbose {
            logger.Errorf("解析目标数据失败: %v", err)
        }

        return false
    }

    if err := s.sendRes(); nil != err {
        if s.verbose {
            logger.Error("发送协议响应数据失败")
        }

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

    if Ver01 != ver[0] {
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
        ipString = string(ip)
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
    data[0] = Ver01
    data[1] = Success
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
