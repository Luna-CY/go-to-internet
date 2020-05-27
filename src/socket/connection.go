package socket

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "net"
    "time"
)

// Connection 隧道连接结构
type Connection struct {
    IsRunning bool
    IsClosed  bool
    Tunnel    net.Conn

    ctx    context.Context
    cancel context.CancelFunc
}

// Init 初始化隧道
func (c *Connection) Init() error {
    c.IsRunning = true
    c.ctx, c.cancel = context.WithCancel(context.Background())

    if _, err := c.Tunnel.Write(make([]byte, 0)); nil != err {
        return err
    }

    return nil
}

// bind 连接隧道
func (c *Connection) Connect(src net.Conn, ipType byte, dstIp string, dstPort int) error {
    connect := tunnel.NewConnectMessage(c.Tunnel, ipType, dstIp, dstPort)
    if err := connect.Send(); nil != err {
        return errors.New(fmt.Sprintf("发送连接消息失败: %v", err))
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
                c.cancel()
                c.Close()

                return err
            }

            over += 1
        case err := <-ch2:
            if nil != err {
                c.cancel()
                c.Close()

                return err
            }

            over += 1
        default:
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
        res, message := tunnel.CopyWithCtxFromMessageProtocol(c.ctx, reader, writer)

        for {
            select {
            case err := <-res:
                if nil != err {
                    c.sendOverMessage()
                    ch <- err

                    return
                }
            case msg := <-message:
                if tunnel.CmdOver != msg.Cmd {
                    c.sendOverMessage()
                    ch <- errors.New(fmt.Sprintf("不支持的消息指令: %v", msg.Cmd))

                    return
                }

                return
            case <-c.ctx.Done():
                c.sendOverMessage()

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
        res := tunnel.CopyLimiterWithCtxToMessageProtocol(c.ctx, reader, writer, nil)

        timer := time.NewTimer(1 * time.Second)
        for {
            select {
            case err := <-res:
                if nil != err {
                    timer.Stop()
                    c.sendOverMessage()

                    ch <- err

                    return
                }

                timer.Reset(1 * time.Second)
            case <-timer.C:
                timer.Stop()
                c.sendOverMessage()

                ch <- nil

                return
            case <-c.ctx.Done():
                timer.Stop()
                c.sendOverMessage()

                return
            }
        }
    }()

    return ch
}

// sendOverMessage 发送结束消息
func (c *Connection) sendOverMessage() {
    if err := tunnel.NewOverMessage(c.Tunnel).Send(); nil != err {
        logger.Errorf("发送结束消息失败: %v", err)
    }
}

// Reset 重置隧道
func (c *Connection) Reset() {
    c.IsRunning = false
    c.ctx = nil
}

// Close 关闭隧道
func (c *Connection) Close() {
    c.Tunnel.Close()
    c.IsClosed = true
}
