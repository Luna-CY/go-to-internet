package proxy

import (
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// StartConnection 开始一个连接处理
func StartConnection(src net.Conn, verbose bool) {
    server, err := tunnel.NewServer(src, verbose)
    if nil != err {
        if verbose {
            logger.Debugf("创建隧道服务端失败: %v", err)
        }

        //ns := http.MockNginx{Conn: &src, Version: "1.14.2"}
        //ns.P404()
        fmt.Println(err)

        return
    }

    if err := server.Bind(); nil != err && verbose {
        logger.Errorf("绑定隧道失败: %v", err)
    }
}
