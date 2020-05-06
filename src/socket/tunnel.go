package socket

import (
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// startTunnel 启动一个隧道
func (s *Socket) startTunnel(src net.Conn, ipType byte, ip string, port int) {
    config := &tunnel.Config{
        ServerHostname: s.Hostname,
        ServerPort:     s.Port,
        ServerUsername: s.Username,
        ServerPassword: s.Password,
        TargetType:     ipType,
        TargetHostOrIp: ip,
        TargetPort:     port,
    }

    dst, err := tunnel.NewClient(config)
    if nil != err {
        fmt.Println("启动隧道失败: ", err)

        return
    }

    dst.Bind(src)
}
