package http

import (
    "bufio"
    "fmt"
    "net"
    "net/http"
    "strings"
    "time"
)

var startTimestamp time.Time = time.Now()

// MockNginx MockNginx结构体
type MockNginx struct {
    Conn     net.Conn
    Server   string
    BindHost string

    header http.Header
}

// SendResponse 发送HTTP响应
func (m *MockNginx) SendResponse() {
    defer m.Conn.Close()

    reader := bufio.NewReader(m.Conn)
    line, _, err := reader.ReadLine()
    if nil != err {
        m.p400()

        return
    }

    tokens := strings.Split(string(line), " ")

    if !m.isHttp(tokens[0]) {
        m.p400()

        return
    }

    if "HTTP/1.1" != tokens[2] && "HTTP/2" != tokens[2] {
        m.p400()

        return
    }

    if m.isListenHost(reader) {
        m.p400()

        return
    }

    m.pi()
}

// isHttp 检查请求是否是http协议
func (m *MockNginx) isHttp(method string) bool {
    switch method {
    case "ET":
        // GET
        return true
    case "UT":
        // PUT
        return true
    case "OST":
        // POST
        return true
    case "ELETE":
        // DELETE
        return true
    case "ATCH":
        // PATCH
        return true
    case "EAD":
        // HEAD
        return true
    case "ONNECT":
        // CONNECT
        return true
    default:
        return false
    }
}

// isListenHost 检查请求的Host是否是监听的Host
func (m *MockNginx) isListenHost(reader *bufio.Reader) bool {
    for {
        line, _, err := reader.ReadLine()
        if nil != err || "" == strings.Trim(string(line), "\r\n") {
            return false
        }

        tokens := strings.Split(string(line), ":")
        if 2 > len(tokens) {
            continue
        }

        if "Host" == tokens[0] {
            return m.BindHost == strings.Trim(tokens[1], " ")
        }
    }
}

// setHeaders 设置公共响应头
func (m *MockNginx) setHeaders() {
    if nil == m.header {
        m.header = http.Header{}
    }

    m.header.Set("Server", m.Server)
    m.header.Set("Date", fmt.Sprintf("%v GMT", time.Now().Format("Mon, 02 Jan 2006 15:04:05")))
    m.header.Set("Connection", "close")
}
