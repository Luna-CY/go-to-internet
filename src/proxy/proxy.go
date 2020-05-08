package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// StartConnection 开始一个连接处理
func StartConnection(src net.Conn, serverConfig *Config, userConfig *config.UserConfig) {
    server, err := tunnel.NewServer(src, userConfig, serverConfig.Verbose)
    if nil != err {
        if serverConfig.Verbose {
            logger.Debugf("创建隧道服务端失败: %v", err)
        }

        ns := http.MockNginx{Conn: src, Server: "nginx", BindHost: serverConfig.Hostname}
        ns.SendResponse()

        return
    }

    if err := server.Bind(); nil != err && serverConfig.Verbose {
        logger.Errorf("绑定隧道失败: %v", err)
    }
}
