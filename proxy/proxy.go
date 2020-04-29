package proxy

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

type Proxy struct{}

// ServeHTTP http请求处理器
func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    fmt.Printf("req from %v to %v\n", request.RemoteAddr, request.URL)
    response, err := p.get(request.URL, &request.Header)
    if nil != err {
        writer.WriteHeader(500)

        _, _ = fmt.Fprintln(writer, "hello word")

        return
    }
    defer response.Body.Close()

    for name, value := range response.Header {
        writer.Header().Add(name, value[0])
    }

    data, _ := ioutil.ReadAll(response.Body)
    writer.WriteHeader(response.StatusCode)

    fmt.Printf("res to %v\n", request.RemoteAddr)
    _, _ = fmt.Fprintf(writer, "%v", string(data))
}

// get http get方法代理
func (p *Proxy) get(url *url.URL, header *http.Header) (*http.Response, error) {
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
