package tunnel

import (
    "errors"
    "fmt"
    "net"
)

// IsTunnelProtocol 检查是否是一个tunnel协议
func IsTunnelProtocol(conn net.Conn) bool {
    return true
}

// ReceiveTarget 获取目标服务器信息
func ReceiveTarget(conn net.Conn) (string, int, error) {
    ver := make([]byte, 1)
    n, err := conn.Read(ver)
    if n != 1 || nil != err {
        fmt.Println(ver)
        return "", 0, errors.New("读取版本号失败")
    }

    port := make([]byte, 2)
    n, err = conn.Read(port)
    if n != 2 || nil != err {
        return "", 0, errors.New("读取端口号失败")
    }

    ipLen := make([]byte, 1)
    n, err = conn.Read(ipLen)
    if n != 1 || nil != err {
        return "", 0, errors.New("读取ip长度失败")
    }

    ip := make([]byte, ipLen[0])
    n, err = conn.Read(ip)
    if n != int(ipLen[0]) || nil != err {
        return "", 0, errors.New("读取ip失败")
    }

    err = sendRes(conn)
    if nil != err {
        return "", 0, err
    }

    return string(ip), int(port[0])<<8 | int(port[1]), nil
}

// sendRes 发送响应数据
func sendRes(conn net.Conn) error {
    dataLength := 1 + 1 + 1 + len("OK")
    data := make([]byte, dataLength)
    data[0] = VER01
    data[1] = 0x00
    data[2] = byte(len("OK"))

    index := 3
    for _, d := range []byte("OK") {
        data[index] = d
        index++
    }

    n, err := conn.Write(data)
    if n != dataLength || nil != err {
        conn.Close()
        return errors.New("写入数据失败")
    }

    return nil
}
