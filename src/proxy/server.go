package proxy

import (
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "io"
    "net"
    "net/http"
    "time"
)

// Server 结构体
type Server struct{}

// ServeHTTP http请求处理器
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    // 不是POST请求暂时忽略
    if "POST" != request.Method {
        return
    }

    reqData, _ := common.ReadAll(request.Body)
    request.Body.Close()

    httpRequest := common.HttpRequest{}
    err := json.Unmarshal(reqData, &httpRequest)
    if nil != err {
        fmt.Printf("解析数据失败: %v\n", err)

        return
    }

    fmt.Printf("%v 请求 %v:%d\n", request.RemoteAddr, httpRequest.TargetIp, httpRequest.TargetPort)

    tcp, err := net.Dial("tcp", fmt.Sprintf("%v:%d", httpRequest.TargetIp, httpRequest.TargetPort))
    if nil != err {
        fmt.Printf("连接目标服务器失败 %v:%d\n", httpRequest.TargetIp, httpRequest.TargetPort)

        return
    }
    defer tcp.Close()

    if _, err = tcp.Write(httpRequest.Data); nil != err {
        fmt.Printf("向目标服务器发送数据失败: %v\n", err)

        return
    }

    _ = tcp.SetReadDeadline(time.Now().Add(3000 * time.Millisecond))
    _, _ = io.Copy(writer, tcp)
}
