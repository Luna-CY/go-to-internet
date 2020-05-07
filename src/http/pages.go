package http

import (
    "compress/gzip"
    "fmt"
)

// p400 400 page
func (m *MockNginx) p400() {
    m.setHeaders()
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

// pi index page
func (m *MockNginx) pi() {
    m.setHeaders()
    m.header.Set("Content-Type", "text/html")
    m.header.Set("Content-Encoding", "gzip")

    content := "<!DOCTYPE html>\n" +
        "<html>\n" +
        "<head>\n" +
        "<title>Welcome to nginx!</title>\n" +
        "<style>\n    body {\n        width: 35em;\n        margin: 0 auto;\n        font-family: Tahoma, Verdana, Arial, sans-serif;\n    }\n</style>\n" +
        "</head>\n" +
        "<body>\n" +
        "<h1>Welcome to nginx!</h1>\n" +
        "<p>If you see this page, the nginx web server is successfully installed and\nworking. Further configuration is required.</p>\n\n" +
        "<p>For online documentation and support please refer to\n<a href=\"http://nginx.org/\">nginx.org</a>.<br/>\nCommercial support is available at\n<a href=\"http://nginx.com/\">nginx.com</a>.</p>\n\n" +
        "<p><em>Thank you for using nginx.</em></p>\n" +
        "</body>\n" +
        "</html>\n"
    m.header.Set("Content-Length", fmt.Sprintf("%d", len(content)))

    // 响应头
    _, _ = m.Conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
    for key, value := range m.header {
        _, _ = m.Conn.Write([]byte(fmt.Sprintf("%v: %v\r\n", key, value[0])))
    }

    // 响应数据
    _, _ = m.Conn.Write([]byte("\r\n"))

    gw := gzip.NewWriter(m.Conn)
    defer gw.Close()
    _, _ = gw.Write([]byte(content))
}
