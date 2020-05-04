package socket

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "net"
    "net/http"
)

var connections map[string]*connection

// init
func init() {
    if nil == connections {
        connections = make(map[string]*connection)
    }
}

// HTTP 结构体
type HTTP struct {
    Sock     *Socket   // Socket结构体指针
    SockConn *net.Conn // socket连接指针

    TargetIp   string // 请求目标的ip
    TargetPort int    // 请求目标的端口
}

// connection 获取连接句柄
func (h *HTTP) connection(ip string, port int) *connection {
    key := fmt.Sprintf("%v:%d", ip, port)
    if conn, ok := connections[key]; ok && !conn.closed {
        return conn
    }

    fmt.Printf("创建到代理服务器的连接: %v:%d\n", ip, port)
    connections[key] = &connection{client: &http.Client{}}

    return connections[key]
}

// request 处理http请求
func (h *HTTP) request() {
    conn := h.connection(h.TargetIp, h.TargetPort)

    reqData, _ := common.ReadAll(*h.SockConn)
    request := common.HttpRequest{TargetIp: h.TargetIp, TargetPort: h.TargetPort, Data: reqData}
    if h.Sock.Authorize {
        request.Username = h.Sock.Username
        request.Password = h.Sock.Password
    }
    body, _ := json.Marshal(request)

    req, _ := http.NewRequest("POST", fmt.Sprintf("https://%v:%d", h.Sock.Hostname, h.Sock.Port), bytes.NewBuffer(body))
    defer req.Body.Close()

    if res, err := conn.do(req); nil == err {
        _, err = common.Copy(*h.SockConn, res.Body)
        if nil != err {
            fmt.Println(err)
        }
    }
}
