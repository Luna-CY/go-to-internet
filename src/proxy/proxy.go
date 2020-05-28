package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "net"
    "runtime"
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
    defer client.Close()

    id := utils.RandomString(8)
    connection := &Connection{Id: id, Tunnel: client, Verbose: p.Verbose}
    if !connection.check(p.UserConfig) {
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
        logger.Infof("建立新的隧道(%v -> %v:%v). ID: %v", connection.Username, connection.UserInfo.CurrentConnection, connection.UserInfo.MaxConnection, id)
    }

    connection.Accept()
    connection.UserInfo.CurrentConnection -= 1

    runtime.GC()
}
