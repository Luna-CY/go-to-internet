package socket

import "net/http"

// connection HTTP连接结构体
type connection struct {
    client *http.Client
    closed bool
}

// close 关闭连接
func (c *connection) close() {
    if nil != c.client {
        c.client.CloseIdleConnections()
    }

    c.closed = true
}

// do 代理client.Do方法
func (c *connection) do(request *http.Request) (*http.Response, error) {
    return c.client.Do(request)
}
