package proxy

import (
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

// StartConnection 开始一个连接处理
func StartConnection(src net.Conn) {
    server, err := tunnel.NewServer(src)
    if nil != err {
        //ns := http.MockNginx{Conn: &src, Version: "1.14.2"}
        //ns.P404()
        fmt.Println(err)

        return
    }

    fmt.Sprintln(server.Bind())
}
