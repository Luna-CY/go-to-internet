package proxy

import (
    "compress/gzip"
    "fmt"
    "net/http"
    "strings"
)

// html 渲染html页面
func (s *Server) html(code int, title string) {
    html := "<html>\r\n<head><title>%v</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>%v</h1></center>\r\n<hr><center>%v</center>\r\n</body>\r\n</html>\r\n"

    s.writer.Header().Set("Server", s.NginxVersion)
    s.writer.Header().Set("Content-Type", "text/html")
    s.writer.Header().Set("Connection", "keep-alive")
    s.writer.Header().Set("Content-Encoding", "gzip")

    gw := gzip.NewWriter(s.writer)
    defer gw.Close()

    flusher, ok := s.writer.(http.Flusher)
    if !ok {
        s.writer.WriteHeader(500)

        title = "500 Internal Server Error"
        _, _ = gw.Write([]byte(fmt.Sprintf(html, title, title, s.NginxVersion)))

        return
    }

    s.writer.WriteHeader(code)

    _, _ = gw.Write([]byte(fmt.Sprintf(html, title, title, s.NginxVersion)))

    if strings.Contains(s.request.Header.Get("User-Agent"), "Chrome") {
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
        _, _ = gw.Write([]byte("<!-- a padding to disable MSIE and Chrome friendly error page -->\r\n"))
    }

    flusher.Flush()
}

// p404 HTTP 404错误页面
func (s *Server) p404() {
    s.html(404, "404 Not Found")
}

// p500 HTTP 500错误页面
func (s *Server) p500() {
    s.html(500, "500 Internal Server Error")
}
