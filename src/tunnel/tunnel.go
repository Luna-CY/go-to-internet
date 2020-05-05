package tunnel

import "net"

// IsTunnelProtocol 检查是否是一个tunnel协议
func IsTunnelProtocol(conn *net.Conn) bool {
    return true
}

func GetTarget(conn *net.Conn) (string, int, error) {
    return "", 0, nil
}
