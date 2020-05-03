package proxy

import (
    "compress/gzip"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "io"
    "net"
    "net/http"
    "time"
)

// Server 结构体
type Server struct {
    NginxVersion string

    request *http.Request
    writer  http.ResponseWriter
}

// ServeHTTP http请求处理器
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    s.request = request
    s.writer = writer

    switch true {
    case "GET" == request.Method && "/" == s.request.RequestURI:
        s.get()
    case "POST" == request.Method && "/" == request.RequestURI:
        s.post()
    default:
        s.p404()
    }
}

// get 处理HTTP GET请求
func (s *Server) get() {
    s.writer.Header().Set("Server", s.NginxVersion)
    s.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
    s.writer.Header().Set("Cache-Control", "no-cache")
    s.writer.Header().Set("Connection", "keep-alive")
    s.writer.Header().Set("Content-Encoding", "gzip")

    response := common.HttpResponse{Code: 0, Message: "OK"}
    body, _ := json.Marshal(response)

    gw := gzip.NewWriter(s.writer)
    defer gw.Close()

    flusher, ok := s.writer.(http.Flusher)
    if !ok {
        s.p500()

        return
    }

    s.writer.WriteHeader(200)
    _, _ = gw.Write(body)

    flusher.Flush()
}

// post 处理HTTP POST请求
func (s *Server) post() {
    reqData, _ := common.ReadAll(s.request.Body)
    s.request.Body.Close()

    httpRequest := common.HttpRequest{}
    err := json.Unmarshal(reqData, &httpRequest)

    // 检查连接的用户身份
    if nil != err || !s.auth(httpRequest.Username, httpRequest.Password) {
        s.p404()

        return
    }

    fmt.Printf("%v 请求 %v:%d\n", s.request.RemoteAddr, httpRequest.TargetIp, httpRequest.TargetPort)

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
    _, _ = io.Copy(s.writer, tcp)
}

// auth 用户鉴权
func (s *Server) auth(username, password string) bool {
    return !("root" != username || "123456" != password)
}
