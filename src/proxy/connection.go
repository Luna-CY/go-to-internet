package proxy

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/time/rate"
    "net"
    "time"
)

// Connection 服务器的连接结构
type Connection struct {
    Id        string
    IsRunning bool
    Limiter   *rate.Limiter
    Tunnel    net.Conn

    Username string
    UserInfo *config.UserInfo
    Protocol *tunnel.HandshakeProtocol

    Verbose bool

    ctx    context.Context
    cancel context.CancelFunc
}

// check 检查连接
func (c *Connection) check(userConfig *config.UserConfig) bool {
    protocol := &tunnel.HandshakeProtocol{Conn: c.Tunnel}
    if err := protocol.ReceiveC(); nil != err {
        return false
    }

    userInfo, ok := userConfig.Users[protocol.Username]
    if !ok {
        return false
    }
    c.Username = protocol.Username
    c.UserInfo = userInfo

    // 检查用户密码
    if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(protocol.Password)); nil != err {
        return false
    }

    // 检查用户有效期
    if "-" != userInfo.Expired {
        expired, err := time.Parse(common.TimePattern, userInfo.Expired)
        if nil != err {
            return false
        }

        if expired.Before(time.Now()) {
            return false
        }
    }
    c.Protocol = protocol

    return true
}

// Accept 接收连接请求并处理
func (c *Connection) Accept() {
    for {
        message := tunnel.NewEmptyMessage(c.Tunnel)
        if err := message.Receive(); nil != err {
            if c.Verbose {
                logger.Errorf("接收消息失败: %v", err)
            }

            return
        }

        if tunnel.CmdNewConnect != message.Cmd {
            if c.Verbose {
                logger.Errorf("不接受的指令: %v", message.Cmd)
            }

            return
        }

        if c.Verbose {
            logger.Infof("新的目标连接请求. 隧道: %v", c.Id)
        }

        if err := message.ParseDst(); nil != err {
            if c.Verbose {
                logger.Errorf("解析目标数据失败: %v", err)
            }

            if err := tunnel.NewOverMessage(c.Tunnel).Send(); nil != err && c.Verbose {
                logger.Errorf("发送结束消息失败: %v", err)
            }

            return
        }

        dst, err := net.Dial("tcp", fmt.Sprintf("%v:%d", message.DstIp, message.DstPort))
        if nil != err {
            if c.Verbose {
                logger.Errorf("建立目标连接失败: %v", err)
            }

            if err := tunnel.NewOverMessage(c.Tunnel).Send(); nil != err && c.Verbose {
                logger.Errorf("发送结束消息失败: %v", err)
            }

            return
        }

        message.Code = tunnel.MessageCodeSuccess
        if err := message.Send(); nil != err {
            if c.Verbose {
                logger.Errorf("发送连接建立响应消息失败: %v", err)
            }

            return
        }

        c.ctx, c.cancel = context.WithCancel(context.Background())
        if err := c.bind(dst); nil != err && c.Verbose {
            logger.Errorf("数据传输失败: %v", err)
        }

        if c.Verbose {
            logger.Info("数据代理完成")
        }

        if err := tunnel.NewOverMessage(c.Tunnel).Send(); nil != err && c.Verbose {
            logger.Errorf("发送结束消息失败: %v", err)
        }
    }
}

// Send 发送消息
func (c *Connection) Send(code byte) error {
    return c.Protocol.Send(code)
}

// bind 连接隧道
func (c *Connection) bind(dst net.Conn) error {
    ch1 := c.bindFromMessage(c.Tunnel, dst)
    ch2 := c.bindToMessage(dst, c.Tunnel)

    over := 0
    for {
        select {
        case err := <-ch1:
            if nil != err {
                c.cancel()

                return err
            }

            over += 1
        case err := <-ch2:
            if nil != err {
                c.cancel()

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
            case <-c.ctx.Done():
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
        var limiter *rate.Limiter
        if 0 != c.UserInfo.MaxRate {
            limiter = rate.NewLimiter(rate.Limit(c.UserInfo.MaxRate*1024), c.UserInfo.MaxRate*512/2)
        }

        res := tunnel.CopyLimiterWithCtxToMessageProtocol(c.ctx, reader, writer, limiter)

        timer := time.NewTimer(1 * time.Second)
        for {
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
                ch <- nil

                return
            }
        }
    }()

    return ch
}
