package socket

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
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
func (c *Connection) Connect(src net.Conn, ipType byte, dstIp string, dstPort int) error {
    connect := tunnel.NewConnectMessage(c.Tunnel, ipType, dstIp, dstPort)
    if err := connect.Send(); nil != err {
        return errors.New(fmt.Sprintf("发送连接消息失败: %v", err))
    }

    if err := connect.Receive(); nil != err {
        return errors.New(fmt.Sprintf("接收连接消息失败: %v", err))
    }

    if tunnel.CmdNewConnect != connect.Cmd || tunnel.MessageCodeSuccess != connect.Code {
        return errors.New(fmt.Sprintf("接收连接消息失败. 响应指令: %v 响应码: %v", connect.Cmd, connect.Code))
    }

    ch1 := c.bindFromMessage(c.Tunnel, src)
    defer close(ch1)
    ch2 := c.bindToMessage(src, c.Tunnel)
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

// bindFromMessage 绑定一个reader和writer
func (c *Connection) bindFromMessage(reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        res := utils.CopyWithCtxFromMessageProtocol(c.ctx, reader, writer)
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

// bindFromMessage 绑定一个reader和writer
func (c *Connection) bindToMessage(reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        res := utils.CopyLimiterWithCtxToMessageProtocol(c.ctx, reader, writer, nil)
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
