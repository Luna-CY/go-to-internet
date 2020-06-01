package socket

import (
    "crypto/tls"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "io"
    "net"
    "runtime"
    "sync"
)

// client 代理客户端结构定义，内部结构
type client struct {
    Socket *Socket

    mutex      sync.Mutex
    maxConnect bool
    stack      *Stack
}

// Init 初始化代理客户端
func (c *client) Init() error {
    c.stack = &Stack{}

    for i := 0; i < 8; i++ {
        conn, err := c.newConnection()
        if nil != err {
            return err
        }

        if nil == conn {
            break
        }

        c.stack.Push(conn)
    }

    return nil
}

// Accept 接收连接请求
func (c *client) Accept(src net.Conn, ipType byte, ip string, port int) {
    defer src.Close()

    for {
        connection, err := c.getConnection()
        if nil != err {
            logger.Errorf("获取隧道失败: %v", err)

            return
        }

        if err := connection.Connect(src, ipType, ip, port); nil != err {
            if ClosedError == err {
                // 如果服务端连接已被关闭，那么重新选择隧道处理此次请求
                continue
            }

            if io.EOF != err {
                connection.Close()
                logger.Errorf("处理请求失败: %v", err)

                // 非隧道关闭异常则处理失败
                return
            }
        }

        if c.Socket.Verbose {
            logger.Info("代理请求完成")
        }

        connection.Reset()
        if !connection.IsClosed {
            c.stack.Push(connection)
        }

        // 该循环在正常情况只执行一次
        break
    }

    runtime.GC()
}

// getConnection 获取一个可用的隧道连接
func (c *client) getConnection() (*Connection, error) {
    for {
        if c.stack.IsEmpty() {
            conn, err := c.newConnection()
            if nil != err {

                return nil, err
            }

            if nil != conn {
                if err := conn.Init(); nil == err {

                    return conn, nil
                }
                conn.Close()

                if c.Socket.Verbose {
                    logger.Errorf("初始化连接失败: %v", err)
                }
            }
        } else {
            conn := c.stack.Pop()

            return conn, nil
        }
    }
}

// newConnection 新建一个隧道连接
func (c *client) newConnection() (*Connection, error) {
    conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%d", c.Socket.Hostname, c.Socket.Port), nil)
    if nil != err {
        return nil, err
    }

    handshake := &tunnel.HandshakeProtocol{Conn: conn, Username: c.Socket.Username, Password: c.Socket.Password}
    if err := handshake.Connection(); nil != err {
        return nil, err
    }

    if err := handshake.ReceiveR(); nil != err {
        return nil, err
    }

    if tunnel.HandshakeCodeConnectionUpperLimit == handshake.Code {
        return nil, nil
    }

    if tunnel.HandshakeCodeSuccess != handshake.Code {
        return nil, errors.New("建立隧道连接失败")
    }

    id := utils.RandomString(8)

    if c.Socket.Verbose {
        logger.Infof("建立新的隧道. ID: %v", id)
    }

    return &Connection{Id: id, Tunnel: conn, Verbose: c.Socket.Verbose}, nil
}
