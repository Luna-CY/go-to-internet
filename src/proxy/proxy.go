package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

type Proxy struct {
    UserConfig *config.UserConfig
    Hostname   string
    Verbose    bool
}

// Init 初始化代理
func (p *Proxy) Init() error {
    return nil
}

// Accept 接收连接请求
func (p *Proxy) Accept(client net.Conn) {
    connection := &Connection{Tunnel: client}
    if !connection.check(p.UserConfig) {
        defer client.Close()

        ns := http.MockNginx{Conn: client, Server: "nginx", BindHost: p.Hostname}
        ns.SendResponse()

        return
    }

    if 0 != connection.UserInfo.MaxConnection && connection.UserInfo.CurrentConnection >= connection.UserInfo.MaxConnection {
        if err := connection.Send(tunnel.HandshakeCodeConnectionUpperLimit); nil != err {
            logger.Errorf("发送握手响应失败: %v", err)
        }

        return
    }

    if err := connection.Send(tunnel.HandshakeCodeSuccess); nil != err {
        logger.Errorf("发送握手响应失败: %v", err)

        return
    }

    connection.UserInfo.CurrentConnection += 1
    if p.Verbose {
        logger.Infof("建立新的连接(%v): 当前 %v 上限 %v", connection.Username, connection.UserInfo.CurrentConnection, connection.UserInfo.MaxConnection)
    }

    connection.Accept()
}
