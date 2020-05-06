package socket

import (
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "io"
    "net"
)

// startTunnel 启动一个隧道
func (s *Socket) startTunnel(src net.Conn, ipType byte, ip string, port int) {
    dst, err := tunnel.StartTunnel(s.Hostname, s.Port, ipType, ip, port)
    if nil != err {
        return
    }

    go func() {
        defer src.Close()
        _, _ = io.Copy(src, dst)

    }()
    _, _ = io.Copy(dst, src)
}
