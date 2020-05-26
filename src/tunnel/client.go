package tunnel

import (
    "crypto/tls"
    "encoding/binary"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "io"
    "net"
    "time"
)

// NewClient 创建一个客户端
func NewClient(config *Config) (*Client, error) {
    conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%d", config.ServerHostname, config.ServerPort), nil)
    if nil != err {
        return nil, err
    }

    client := &Client{serverConn: conn, config: config}
    if err = client.connect(); nil != err {
        defer conn.Close()
        return nil, err
    }

    return client, nil
}

// CheckServer 检查服务端是否支持隧道协议
func CheckServer(config *Config) error {
    conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%d", config.ServerHostname, config.ServerPort), nil)
    if nil != err {
        return err
    }
    defer conn.Close()

    config.TargetType = 0x01
    config.TargetHostOrIp = CheckConnectTargetIp
    config.TargetPort = CheckConnectTargetPort

    client := &Client{serverConn: conn, config: config}
    return client.connect()
}

// client 隧道的客户端结构体
type Client struct {
    serverConn net.Conn // 服务器连接
    config     *Config
}

// bind 双向绑定服务端以及请求来源
func (c *Client) Bind(src net.Conn) error {
    defer src.Close()

    go func() {
        defer c.serverConn.Close()
        _, _ = io.Copy(c.serverConn, src)
    }()
    _, _ = io.Copy(src, c.serverConn)

    return nil
}

// connect 连接服务器
func (c *Client) connect() error {
    for {
        if err := c.sendConnectData(); nil != err {
            return err
        }

        code, msg, err := c.receive()
        if nil != err {
            return err
        }

        if CodeSuccess == code {
            break
        } else if CodeConnectionUpperLimit == code {
            time.Sleep(1 * time.Second)

            continue
        } else {
            logger.Errorf("无法识别的消息: %v", msg)
        }
    }

    return nil
}

// sendConnectData 发送用户信息
func (c *Client) sendConnectData() error {
    c.config.ServerPassword = utils.EncryptPassword(c.config.ServerPassword)

    userInfoLen := 1 + 1 + len(c.config.ServerUsername) + 1 + len(c.config.ServerPassword)
    targetInfoLen := 2 + 1 + 1 + len(c.config.TargetHostOrIp)
    dataLength := userInfoLen + targetInfoLen

    data := make([]byte, dataLength)
    data[0] = HandshakeProtocolVersion
    data[1] = byte(len(c.config.ServerUsername))

    index := 2
    for _, d := range []byte(c.config.ServerUsername) {
        data[index] = d
        index++
    }

    data[index] = byte(len(c.config.ServerPassword))
    index++

    for _, d := range []byte(c.config.ServerPassword) {
        data[index] = d
        index++
    }

    bs := make([]byte, 2)
    binary.BigEndian.PutUint16(bs, uint16(c.config.TargetPort))

    for _, d := range bs {
        data[index] = d
        index++
    }

    data[index] = c.config.TargetType
    index++

    data[index] = byte(len(c.config.TargetHostOrIp))
    index++

    for _, d := range []byte(c.config.TargetHostOrIp) {
        data[index] = d
        index++
    }

    n, err := c.serverConn.Write(data)
    if n != dataLength || nil != err {
        c.serverConn.Close()

        return errors.New("写入数据失败")
    }

    return nil
}

// receiveRes 读取响应消息
func (c *Client) receive() (byte, string, error) {
    ver := make([]byte, 1)
    n, err := c.serverConn.Read(ver)
    if n != 1 || nil != err {
        c.serverConn.Close()

        return 0xff, "", errors.New("读取应答版本号失败")
    }

    if HandshakeProtocolVersion != ver[0] {
        return 0xff, "", errors.New("不支持的协议版本")
    }

    code := make([]byte, 1)
    n, err = c.serverConn.Read(code)
    if n != 1 || nil != err {
        c.serverConn.Close()

        return 0xff, "", errors.New("读取响应码失败")
    }

    msgLen := make([]byte, 1)
    n, err = c.serverConn.Read(msgLen)
    if n != 1 || nil != err {
        c.serverConn.Close()

        return 0xff, "", errors.New("读取消息长度失败")
    }

    msg := make([]byte, msgLen[0])
    n, err = c.serverConn.Read(msg)
    if n != int(msgLen[0]) || nil != err {
        c.serverConn.Close()

        return 0xff, "", errors.New("读取消息失败")
    }

    return code[0], string(msg), nil
}
