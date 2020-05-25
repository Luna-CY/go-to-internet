package tunnel

import (
    "context"
    "net"
    "sync"
)

var servers []*Connection
var clients []*Connection

func GetConnection() (*Connection, error) {
    return nil, nil
}

type Connection struct {
    ctx context.Context
    mutex sync.Mutex
    tunnel net.Conn
}
