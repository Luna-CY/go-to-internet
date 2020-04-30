package proxy

import (
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "net"
    "net/http"
)

// Server 结构体
type Server struct {
    responseWriter http.ResponseWriter
}

// ServeHTTP http请求处理器
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    s.responseWriter = writer

    // 不是POST请求暂时忽略
    if "POST" != request.Method {
        return
    }

    httpRequest := common.HttpRequest{}

    reqData, _ := common.ReadAll(request.Body)
    request.Body.Close()

    err := json.Unmarshal(reqData, &httpRequest)
    if nil != err {
        s.send(common.UnSerializeDataFail, "解析请求数据失败", nil, true)

        return
    }

    fmt.Printf("%v 请求 %v:%d\n", request.RemoteAddr, httpRequest.TargetIp, httpRequest.TargetPort)

    tcp, err := net.Dial("tcp", fmt.Sprintf("%v:%d", httpRequest.TargetIp, httpRequest.TargetPort))
    if nil != err {
        fmt.Printf("连接目标服务器失败 %v:%d\n", httpRequest.TargetIp, httpRequest.TargetPort)
        s.send(common.ConnectToTargetFail, "代理失败: 无法连接目标服务器", nil, true)

        return
    }
    defer tcp.Close()

    if _, err = tcp.Write(httpRequest.Data); nil != err {
        s.send(common.WriteDataToTargetFail, "代理失败: 向目标服务器发送数据失败", nil, true)

        return
    }

    for {
        isLast := false
        resData, _ := common.ReadAll(tcp)

        if "\r\n\r\n" == string(resData[len(resData)-4:]) {
            isLast = true
        }

        s.send(common.Success, "", resData, isLast)
        if isLast {
            break
        }
    }
}

// send 向http请求发送响应
func (s *Server) send(code int, message string, data []byte, isLast bool) {
    if nil == s.responseWriter {
        return
    }

    httpResponse := common.HttpResponse{Code: code, Message: message, Data: data, IsLast: isLast}
    body, _ := json.Marshal(httpResponse)

    _, _ = s.responseWriter.Write(body)
}
