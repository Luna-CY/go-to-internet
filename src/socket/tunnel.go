package socket

import (
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// startTunnel 启动一个隧道
func (s *Socket) startTunnel(src net.Conn, ipType byte, ip string, port int, verbose bool) {
    config := &tunnel.Config{
        ServerHostname: s.Hostname,
        ServerPort:     s.Port,
        ServerUsername: s.Username,
        ServerPassword: s.Password,
        TargetType:     ipType,
        TargetHostOrIp: ip,
        TargetPort:     port,
        Verbose:        verbose,
    }

    dst, err := tunnel.NewClient(config)
    if nil != err {
        if verbose {
            logger.Errorf("启动隧道失败: %v", err)
        }

        return
    }

    if err := dst.Bind(src); nil != err && verbose {
        logger.Errorf("绑定隧道失败: %v", err)
    }
}
