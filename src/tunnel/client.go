package tunnel

import (
    "crypto/tls"
    "encoding/binary"
    "errors"
    "fmt"
    "net"
)

func StartTunnel(serverHost string, serverPort int, ipType byte, ip string, port int) (net.Conn, error) {
    connection, err := tls.Dial("tcp", fmt.Sprintf("%v:%d", serverHost, serverPort), nil)
    if nil != err {
        return nil, err
    }

    if err = sendTarget(connection, ipType, ip, port); nil != err {
        return nil, err
    }

    if err = receiveRes(connection); nil != err {
        return nil, err
    }

    return connection, nil
}

// sendTarget 发送target信息
func sendTarget(connection net.Conn, ipType byte, ip string, port int) error {
    dataLength := 1 + 2 + 1 + 1 + len(ip)
    data := make([]byte, dataLength)
    data[0] = VER01

    bs := make([]byte, 2)
    binary.BigEndian.PutUint16(bs, uint16(port))

    index := 1
    for _, d := range bs {
        data[index] = d
        index++
    }

    data[index] = ipType
    index++

    data[index] = byte(len(ip))
    index++

    for _, d := range []byte(ip) {
        data[index] = d
        index++
    }

    fmt.Println("发送目标信息: ", data, " -> ", string(data))
    n, err := connection.Write(data)
    if n != dataLength || nil != err {
        connection.Close()
        return errors.New("写入数据失败")
    }

    return nil
}

// receiveRes 读取响应消息
func receiveRes(connection net.Conn) error {
    ver := make([]byte, 1)
    n, err := (connection).Read(ver)
    if n != 1 || nil != err {
        return errors.New("读取应答版本号失败")
    }

    code := make([]byte, 1)
    n, err = (connection).Read(code)
    if n != 1 || nil != err {
        return errors.New("读取响应码失败")
    }

    msgLen := make([]byte, 1)
    n, err = (connection).Read(msgLen)
    if n != 1 || nil != err {
        return errors.New("读取消息长度失败")
    }

    msg := make([]byte, msgLen[0])
    n, err = (connection).Read(msg)
    if n != int(msgLen[0]) || nil != err {
        return errors.New("读取消息失败")
    }

    return nil
}
