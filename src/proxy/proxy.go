package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// StartConnection 开始一个连接处理
func StartConnection(src net.Conn, hostname string, verbose bool) {
    server, err := tunnel.NewServer(src, verbose)
    if nil != err {
        if verbose {
            logger.Debugf("创建隧道服务端失败: %v", err)
        }

        ns := http.MockNginx{Conn: src, Server: "nginx", BindHost: hostname}
        ns.SendResponse()

        return
    }

    if err := server.Bind(); nil != err && verbose {
        logger.Errorf("绑定隧道失败: %v", err)
    }
}
