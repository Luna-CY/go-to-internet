package socket

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/common"
    "golang.org/x/net/http2"
    "io"
    "net"
    "net/http"
)

// HTTP 结构体
type HTTP struct {
    Sock *Socket   // Socket结构体指针
    Conn *net.Conn // socket连接指针

    TargetIp   string // 请求目标的ip
    TargetPort int    // 请求目标的端口
}

// request 处理http请求
func (h *HTTP) request() {
    // 创建http/2客户端
    client := http.Client{
        Transport: &http2.Transport{
            AllowHTTP: true,
            DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
                return net.Dial(network, addr)
            },
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    var db []byte
    data := make([]byte, 256)

    for {
        n, err := (*h.Conn).Read(data)
        db = append(db, data[:n]...)

        if io.EOF == err || n < len(data) {
            break
        }
    }

    request := common.HttpRequest{TargetIp: h.TargetIp, TargetPort: h.TargetPort, Data: db}
    body, _ := json.Marshal(request)
    fmt.Println(string(body))

    fmt.Printf("request to https://%v:%d\n", h.Sock.Hostname, h.Sock.Port)
    req, _ := http.NewRequest("POST", fmt.Sprintf("https://%v:%d", h.Sock.Hostname, h.Sock.Port), bytes.NewBuffer(body))

    res, err := client.Do(req)
    if nil != err {
        fmt.Println(err)

        return
    }

    fmt.Println(res.StatusCode)
}
