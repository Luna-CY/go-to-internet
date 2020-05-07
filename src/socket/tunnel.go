package socket

import (
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
)

// startTunnel 启动一个隧道
func (s *Socket) startTunnel(ipType byte, ip string, port int, verbose bool) (*tunnel.Client, error) {
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

    client, err := tunnel.NewClient(config)
    if nil != err {
        if verbose {
            logger.Errorf("启动隧道失败: %v", err)
        }

        return nil, err
    }

    return client, nil
}
