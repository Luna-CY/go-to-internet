package socket

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
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
    client := http.Client{}

    reqData, _ := common.ReadAll(*h.Conn)
    request := common.HttpRequest{TargetIp: h.TargetIp, TargetPort: h.TargetPort, Data: reqData}
    body, _ := json.Marshal(request)

    fmt.Printf("request to %v:%d\n", request.TargetIp, request.TargetPort)
    req, _ := http.NewRequest("POST", fmt.Sprintf("https://%v:%d", h.Sock.Hostname, h.Sock.Port), bytes.NewBuffer(body))

    res, err := client.Do(req)
    req.Body.Close()

    if nil != err {
        fmt.Println(err)

        return
    }

    httpResponse := common.HttpResponse{}

    for {
        resData, _ := common.ReadAll(res.Body)

        _ = json.Unmarshal(resData, &httpResponse)
        if common.Success != httpResponse.Code {
            fmt.Println(httpResponse.Message)

            return
        }

        _, _ = (*h.Conn).Write(httpResponse.Data)
        if httpResponse.IsLast {
            res.Body.Close()

            break
        }
    }
}
