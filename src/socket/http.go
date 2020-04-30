package socket

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "golang.org/x/net/http2"
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
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    _ = http2.ConfigureTransport(transport)

    client := http.Client{Transport: transport}

    reqData, _ := common.ReadAll(*h.Conn)
    request := common.HttpRequest{TargetIp: h.TargetIp, TargetPort: h.TargetPort, Data: reqData}
    body, _ := json.Marshal(request)

    fmt.Printf("request to https://%v:%d\n", h.Sock.Hostname, h.Sock.Port)
    req, _ := http.NewRequest("POST", fmt.Sprintf("https://%v:%d", h.Sock.Hostname, h.Sock.Port), bytes.NewBuffer(body))

    res, err := client.Do(req)
    if nil != err {
        fmt.Println(err)

        return
    }

    httpResponse := common.HttpResponse{}
    resData, _ := common.ReadAll(res.Body)
    _ = json.Unmarshal(resData, &httpResponse)

    if common.Success != httpResponse.Code {
        fmt.Println(httpResponse.Message)

        return
    }

    _, _ = (*h.Conn).Write(httpResponse.Data)
}
