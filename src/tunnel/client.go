package tunnel

import (
    "crypto/tls"
    "encoding/binary"
    "errors"
    "fmt"
    "net"
)

// NewClient 创建一个客户端
func NewClient(config *Config) (*Client, error) {
    conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%d", config.ServerHostname, config.ServerPort), nil)
    if nil != err {
        return nil, err
    }

    client := &Client{conn: conn}
    if err = client.connect(); nil != err {
        return nil, err
    }

    return client, nil
}

// Client 隧道的客户端结构体
type Client struct {
    conn   net.Conn
    config *Config
}

// Bind 双向绑定服务端以及请求来源
func (c *Client) Bind(src net.Conn) {}

// connect 连接服务器
func (c *Client) connect() error {
    // TODO: 发送用户验证
    if err := c.sendTarget(); nil != err {
        return err
    }

    if err := c.receiveRes(); nil != err {
        return err
    }

    return nil
}

// sendTarget 发送target信息
func (c *Client) sendTarget() error {
    dataLength := 1 + 2 + 1 + 1 + len(c.config.TargetHostOrIp)
    data := make([]byte, dataLength)
    data[0] = VER01

    bs := make([]byte, 2)
    binary.BigEndian.PutUint16(bs, uint16(c.config.TargetPort))

    index := 1
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

    fmt.Println("发送目标信息: ", data, " -> ", string(data))
    n, err := c.conn.Write(data)
    if n != dataLength || nil != err {
        c.conn.Close()

        return errors.New("写入数据失败")
    }

    return nil
}

// receiveRes 读取响应消息
func (c *Client) receiveRes() error {
    ver := make([]byte, 1)
    n, err := c.conn.Read(ver)
    if n != 1 || nil != err {
        c.conn.Close()

        return errors.New("读取应答版本号失败")
    }

    code := make([]byte, 1)
    n, err = c.conn.Read(code)
    if n != 1 || nil != err {
        c.conn.Close()

        return errors.New("读取响应码失败")
    }

    msgLen := make([]byte, 1)
    n, err = c.conn.Read(msgLen)
    if n != 1 || nil != err {
        c.conn.Close()

        return errors.New("读取消息长度失败")
    }

    msg := make([]byte, msgLen[0])
    n, err = c.conn.Read(msg)
    if n != int(msgLen[0]) || nil != err {
        c.conn.Close()

        return errors.New("读取消息失败")
    }

    return nil
}
