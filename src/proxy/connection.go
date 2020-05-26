package proxy

import (
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "golang.org/x/crypto/bcrypt"
    "net"
    "time"
)

// connection 客户端与服务器的连接结构
type connection struct {
    Client net.Conn

    username string
    userInfo *config.UserInfo
    protocol *tunnel.Protocol
}

// check 检查连接
func (c *connection) check(userConfig *config.UserConfig) bool {
    protocol := &tunnel.Protocol{Conn: c.Client}
    if err := protocol.Receive(); nil != err {
        return false
    }

    userInfo, ok := userConfig.Users[protocol.GetUsername()]
    if !ok {
        return false
    }
    c.username = protocol.GetUsername()
    c.userInfo = userInfo

    // 检查用户密码
    if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(protocol.GetPassword())); nil != err {
        return false
    }

    // 检查用户有效期
    if "-" != userInfo.Expired {
        expired, err := time.Parse(common.TimePattern, userInfo.Expired)
        if nil != err {
            return false
        }

        if expired.After(time.Now()) {
            return false
        }
    }
    c.protocol = protocol

    return true
}

// send 发送消息
func (c *connection) send(code byte) error {
    return c.protocol.Send(code)
}
