package tunnel

import (
    "context"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "golang.org/x/time/rate"
    "io"
    "net"
    "sync"
    "time"
)

// Connection 隧道连接结构
type Connection struct {
    IsRunning bool
    Limiter   *rate.Limiter
    Tunnel    net.Conn

    ctx   context.Context
    mutex sync.Mutex
}

// Init 初始化隧道
func (c *Connection) Init() error {
    c.mutex.Lock()
    c.IsRunning = true
    c.ctx = context.Background()

    if _, err := c.Tunnel.Read(make([]byte, 0)); nil != err {
        return err
    }

    return nil
}

// Connect 连接隧道
func (c *Connection) Connect(dst net.Conn) error {
    ch1 := c.bind(c.Tunnel, dst)
    defer close(ch1)
    ch2 := c.bind(dst, c.Tunnel)
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
        res := utils.CopyLimiterWithCtx(c.ctx, reader, writer, c.Limiter)
        defer close(res)

        timer := time.NewTimer(3 * time.Second)
        select {
        case err := <-res:
            if nil != err {
                timer.Stop()

                ch <- err

                return
            }

            timer.Reset(3 * time.Second)
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

// 重置隧道
func (c *Connection) Reset() {
    c.mutex.Unlock()
    c.IsRunning = false
    c.ctx = nil
}
