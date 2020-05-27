package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

type Proxy struct {
    UserConfig *config.UserConfig
    Hostname   string
    Verbose    bool

    connections map[string][]*tunnel.Connection
}

// Init 初始化代理
func (p *Proxy) Init() error {
    return nil
}

// Accept 接收连接请求
func (p *Proxy) Accept(client net.Conn) {
    connection := &Connection{Client: client}
    if !connection.check(p.UserConfig) {
        defer client.Close()

        ns := http.MockNginx{Conn: client, Server: "nginx", BindHost: p.Hostname}
        ns.SendResponse()

        return
    }
}
