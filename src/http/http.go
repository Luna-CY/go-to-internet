package http

import (
    "compress/gzip"
    "fmt"
    "net"
    "net/http"
    "time"
)

// MockNginx MockNginx结构体
type MockNginx struct {
    Conn   net.Conn
    Server string

    header http.Header
}

// SendResponse 发送HTTP响应
func (m *MockNginx) SendResponse() {
    defer m.Conn.Close()

    if nil == m.header {
        m.header = http.Header{}
    }

    m.header.Set("Server", m.Server)
    m.header.Set("Date", time.Now().Format(time.RFC1123))

    if duration, err := time.ParseDuration("-8h"); nil == err {
        m.header.Set("Date", fmt.Sprintf("%v GMT", time.Now().Add(duration).Format("Mon, 02 Jan 2006 15:04:05")))
    }

    m.header.Set("Connection", "keep-alive")
    m.header.Set("Content-Type", "application/json; charset=utf-8")
    m.header.Set("Content-Encoding", "gzip")

    gw := gzip.NewWriter(m.Conn)
    defer gw.Close()

    // 响应头
    _, _ = m.Conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
    for key, value := range m.header {
        _, _ = m.Conn.Write([]byte(fmt.Sprintf("%v: %v\r\n", key, value[0])))
    }

    // 响应数据
    _, _ = m.Conn.Write([]byte("\r\n"))
    _, _ = gw.Write([]byte("{\"code\": 1, \"msg\": \"签名错误\"}"))
}
