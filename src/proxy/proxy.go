package proxy

import (
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/http"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
)

var connections map[string]*net.Conn

// init
func init() {
    if nil == connections {
        connections = make(map[string]*net.Conn)
    }
}

// Close 关闭连接池
func Close() {
    if nil != connections {
        for key := range connections {
            fmt.Println("关闭连接: ", key)
            (*connections[key]).Close()
        }
    }
}

func connection(ip string, port int) *net.Conn {
    key := fmt.Sprintf("%v:%d", ip, port)
    if conn, ok := connections[key]; ok {
        return conn
    }

    fmt.Printf("创建到目标服务器连接: %v:%d\n", ip, port)
    connection, err := net.Dial("tcp", fmt.Sprintf("%v:%d", ip, port))
    if nil != err {
        fmt.Printf("连接目标服务器失败 %v:%d\n", ip, port)

        return nil
    }
    connections[key] = &connection

    return connections[key]
}

// StartConnection 开始一个连接处理
func StartConnection(src net.Conn, config *ServerConfig) {
    if !tunnel.IsTunnelProtocol(src) {
        ns := http.MockNginx{Conn: &src, Version: "1.14.2"}
        ns.P404()

        return
    }

    ip, port, err := tunnel.ReceiveTarget(src)
    if nil != err {
        fmt.Println(err)

        return
    }

    fmt.Println("ip: ", ip, " -> ", port)

    //dst := connection(ip, port)
    //_, _ = io.Copy(*dst, src)
    //_, _ = io.Copy(src, *dst)
}

func CheckUser(username, password string) bool {
    return true
}
