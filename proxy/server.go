package proxy

import (
    "encoding/json"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/common"
    "io"
    "net/http"
    "net/url"
)

// Server 结构体
type Server struct{}

// ServeHTTP http请求处理器
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    fmt.Printf("req from %v to %v\n", request.RemoteAddr, request.URL)

    fmt.Printf("URL: %v\n", request.URL)
    fmt.Printf("METHOD: %v\n", request.Method)

    var db []byte
    data := make([]byte, 256)

    for {
        n, err := request.Body.Read(data)
        db = append(db, data[:n]...)

        if io.EOF == err {
            break
        }
    }

    httpRequest := common.HttpRequest{}
    err := json.Unmarshal(db, &httpRequest)
    if nil != err {
        writer.WriteHeader(400)

        return
    }

    fmt.Printf("TARGET IP: %v\n", httpRequest.TargetIp)
    fmt.Printf("TARGET PORT: %d\n", httpRequest.TargetPort)
    fmt.Printf("DATA: %v\n", httpRequest.Data)

    writer.WriteHeader(200)
    _, _ = fmt.Fprintf(writer, "请求成功")
}

// get http get方法代理
func (s *Server) get(url *url.URL, header *http.Header) (*http.Response, error) {
    request, err := http.NewRequest("GET", url.String(), nil)
    if nil != err {
        return nil, err
    }

    for name, value := range *header {
        request.Header.Add(name, value[0])
    }

    client := &http.Client{}
    response, err := client.Do(request)
    if nil != err {
        return nil, err
    }
    defer client.CloseIdleConnections()

    return response, nil
}
