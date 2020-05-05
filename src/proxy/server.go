package proxy

import (
    "compress/gzip"
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "net"
    "net/http"
)

// Server 结构体
type Server struct {
    NginxVersion string
    Conn         net.Conn

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
    s.request.Body.Close() // 直接关闭通道，不接收任何数据

    s.writer.Header().Set("Server", s.NginxVersion)
    s.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
    s.writer.Header().Set("Cache-Control", "no-cache")
    s.writer.Header().Set("Connection", "keep-alive")
    s.writer.Header().Set("Content-Encoding", "gzip")

    response := common.HttpResponse{Code: 0, Message: "OK"}
    body, _ := json.Marshal(response)

    gw := gzip.NewWriter(s.writer)
    defer gw.Close()

    if flusher, ok := s.writer.(http.Flusher); ok {
        s.writer.WriteHeader(200)
        _, _ = gw.Write(body)

        flusher.Flush()

        return
    }

    s.p500()
}

// connection 获取tcp连接句柄
func (s *Server) connection(ip string, port int) *net.Conn {
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

// post 处理HTTP POST请求
func (s *Server) post() {
    reqData, _ := common.ReadAll(s.request.Body)
    defer s.request.Body.Close()

    httpRequest := common.HttpRequest{}
    err := json.Unmarshal(reqData, &httpRequest)

    // 检查连接的用户身份
    if nil != err || !s.auth(httpRequest.Username, httpRequest.Password) {
        s.p404()

        return
    }

    if conn := s.connection(httpRequest.TargetIp, httpRequest.TargetPort); nil != conn {
        if _, err = (*conn).Write(httpRequest.Data); nil != err {
            fmt.Printf("向目标服务器发送数据失败: %v\n", err)

            return
        }

        flusher, ok := s.writer.(http.Flusher)
        if ok {
            s.writer.Header().Set("Content-Type", "application/octet-stream")
            s.writer.Header().Set("Content-Length", "-1")

            var data []byte
            buffer := make([]byte, 256)
            for {
                n, err := (*conn).Read(buffer)
                _, _ = s.writer.Write(buffer[:n])
                data = append(data, buffer[:n]...)
                flusher.Flush()

                if nil != err || len(buffer) > n {
                    break
                }
            }

            fmt.Println(string(data))
            s.writer.WriteHeader(200)
        }
    }
}

// auth 用户鉴权
func (s *Server) auth(username, password string) bool {
    return !("root" != username || "123456" != password)
}
