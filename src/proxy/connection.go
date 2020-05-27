package proxy

import (
    "context"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/time/rate"
    "net"
    "time"
)

// Connection 服务器的连接结构
type Connection struct {
    IsRunning bool
    Limiter   *rate.Limiter
    Tunnel    net.Conn

    Username string
    UserInfo *config.UserInfo
    Protocol *tunnel.HandshakeProtocol

    ctx context.Context
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

func (c *Connection) Accept() {
    for {
        message := tunnel.NewEmptyMessage(c.Tunnel)
        if err := message.Receive(); nil != err {
            logger.Errorf("接收消息失败: %v", err)

            continue
        }


    }
}

// Send 发送消息
func (c *Connection) Send(code byte) error {
    return c.Protocol.Send(code)
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
func (c *Connection) bind(reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        res := utils.CopyLimiterWithCtxToMessageProtocol(c.ctx, reader, writer, c.Limiter)
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
