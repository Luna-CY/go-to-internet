package http

import (
    "compress/gzip"
    "net"
    "net/http"
)

type MockNginx struct {
    Conn    *net.Conn
    Version string

    header http.Header
}

func (m *MockNginx) sendResponse(code int, content string) {
    m.header.Set("Server", m.Version)
    m.header.Set("Connection", "keep-alive")
    m.header.Set("Content-Type", "text/html")
    m.header.Set("Content-Encoding", "gzip")

    gw := gzip.NewWriter(*m.Conn)
    defer gw.Close()
}
