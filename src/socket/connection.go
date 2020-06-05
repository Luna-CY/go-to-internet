package socket

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "math/rand"
    "net"
    "time"
)

// Connection 隧道连接结构
type Connection struct {
    Id        string
    IsRunning bool
    IsClosed  bool
    IsTimeout bool
    Tunnel    net.Conn

    Verbose bool
}

// Init 初始化隧道
func (c *Connection) Init() error {
    c.IsRunning = true
    if _, err := c.Tunnel.Write(make([]byte, 0)); nil != err {
        return err
    }

    timeout := time.Duration(rand.Intn(3000)+600) * time.Second
    time.AfterFunc(timeout, c.timeout)

    return nil
}

// bind 连接隧道
func (c *Connection) Connect(src net.Conn, ipType byte, dstIp string, dstPort int) error {
    if c.Verbose {
        logger.Infof("发起新的数据连接请求. 隧道: %v", c.Id)
    }

    connect := tunnel.NewConnectMessage(c.Tunnel, ipType, dstIp, dstPort)
    if err := connect.Send(); nil != err {
        return ClosedError
    }

    if err := connect.Receive(); nil != err {
        return errors.New(fmt.Sprintf("接收连接消息失败: %v", err))
    }

    if tunnel.CmdNewConnect != connect.Cmd || tunnel.MessageCodeSuccess != connect.Code {
        return errors.New(fmt.Sprintf("建立连接失败. 响应指令: %v 响应码: %v", connect.Cmd, connect.Code))
    }

    ch1 := c.bindFromMessage(c.Tunnel, src)
    ch2 := c.bindToMessage(src, c.Tunnel)

    over := 0
    for {
        select {
        case err := <-ch1:
            if nil != err {
                c.Close()

                return err
            }

            over += 1
            if over == 2 {
                return nil
            }
        case err := <-ch2:
            if nil != err {
                c.Close()

                return err
            }
            c.sendOverMessage()

            over += 1
            if over == 2 {
                return nil
            }
        }
    }
}

// bindFromMessage 绑定一个reader和writer
func (c *Connection) bindFromMessage(reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        res, message := tunnel.CopyFromMessageProtocol(reader, writer)

        for {
            select {
            case err := <-res:
                if nil != err {
                    ch <- err

                    return
                }
            case msg := <-message:
                if tunnel.CmdOver != msg.Cmd {
                    ch <- errors.New(fmt.Sprintf("不支持的消息指令: %v", msg.Cmd))

                    return
                }
                ch <- nil

                return
            }
        }
    }()

    return ch
}

// bindToMessage 绑定一个reader和writer
func (c *Connection) bindToMessage(reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        ctx, cancel := context.WithCancel(context.Background())
        res := tunnel.CopyLimiterWithCtxToMessageProtocol(ctx, reader, writer, nil)

        timer := time.NewTimer(30 * time.Second)
        for {
            select {
            case err := <-res:
                if nil != err {
                    timer.Stop()
                    ch <- err

                    return
                }

                timer.Reset(30 * time.Second)
            case <-timer.C:
                cancel()
                timer.Stop()
                ch <- nil

                return
            }
        }
    }()

    return ch
}

// sendOverMessage 发送结束消息
func (c *Connection) sendOverMessage() {
    if err := tunnel.NewOverMessage(c.Tunnel).Send(); nil != err && c.Verbose {
        logger.Errorf("发送结束消息失败: %v", err)
    }
}

// timeout 连接已经达到超时时间
func (c *Connection) timeout() {
    c.IsTimeout = true

    if !c.IsRunning {
        c.Close()
    }
}

// Reset 重置隧道
func (c *Connection) Reset() {
    c.IsRunning = false
}

// Close 关闭隧道
func (c *Connection) Close() {
    c.Tunnel.Close()
    c.IsClosed = true
}
