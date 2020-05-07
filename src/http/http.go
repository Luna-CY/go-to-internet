package http

import (
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

    m.header.Set("Connection", "close")
    m.header.Set("Content-Type", "text/html")

    content := "<html>\r\n<head><title>%v</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>%v</h1></center>\r\n<hr><center>%v</center>\r\n</body>\r\n</html>\r\n"
    content = fmt.Sprintf(content, "400 Bad Request", "400 Bad Request", m.Server)
    m.header.Set("Content-Length", fmt.Sprintf("%d", len(content)))

    // 响应头
    _, _ = m.Conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
    for key, value := range m.header {
        _, _ = m.Conn.Write([]byte(fmt.Sprintf("%v: %v\r\n", key, value[0])))
    }

    // 响应数据
    _, _ = m.Conn.Write([]byte("\r\n"))
    _, _ = m.Conn.Write([]byte(content))
}
