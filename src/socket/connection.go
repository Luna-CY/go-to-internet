package socket

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "io"
    "net"
    "sync"
    "time"
)

// Connection 隧道连接结构
type Connection struct {
    IsRunning bool
    Tunnel    net.Conn

    ctx   context.Context
    mutex sync.Mutex
}

// Lock 锁定当前隧道
func (c *Connection) Lock() {
    c.mutex.Lock()
    c.IsRunning = true
}

// Init 初始化隧道
func (c *Connection) Init() error {
    c.ctx = context.Background()

    if _, err := c.Tunnel.Write(make([]byte, 0)); nil != err {
        return err
    }

    return nil
}

// Connect 连接隧道
func (c *Connection) Connect(src net.Conn, dstIp string, dstPort int) error {
    connect := tunnel.NewConnectMessage(c.Tunnel, dstIp, dstPort)
    if err := connect.Send(); nil != err {
        return errors.New(fmt.Sprintf("建立连接失败: %v", err))
    }

    ch1 := c.bind(c.Tunnel, src)
    defer close(ch1)
    ch2 := c.bind(src, c.Tunnel)
    defer close(ch2)

    over := 0

    select {
    case err := <-ch1:
        if nil != err {
            return err
        }
        over += 1
    case err := <-ch2:
        if nil != err {
            return err
        }

        over += 1
    default:
        if over == 2 {
            break
        }
    }

    return nil
}

// bind 绑定一个reader和writer
func (c *Connection) bind(reader io.Reader, writer io.Writer) chan error {
    ch := make(chan error)

    go func() {
        res := utils.CopyWithCtx(c.ctx, reader, writer)
        defer close(res)

        timer := time.NewTimer(1 * time.Second)
        select {
        case err := <-res:
            if nil != err {
                timer.Stop()

                ch <- err

                return
            }

            timer.Reset(1 * time.Second)
        case <-timer.C:
            timer.Stop()

            ch <- nil

            return
        case <-c.ctx.Done():
            timer.Stop()

            return
        }
    }()

    return ch
}

// Close 重置隧道
func (c *Connection) Close() {
    c.mutex.Unlock()
    c.IsRunning = false
    c.ctx = nil
}
