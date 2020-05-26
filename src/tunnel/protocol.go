package tunnel

// 握手协议
//
// 建立连接
// VER USER_L USER PASS_L PASS
//  1    1     N     1     N
//
// 响应消息
// VER CODE
//  1   1
//
// 通信协议
//
// VER CMD DATA_L DATA
//  1   1    1     N

const CodeSuccess = 0x00
const CodeConnectionUpperLimit = 0x01

const CheckConnectTargetIp = "0.0.0.0"
const CheckConnectTargetPort = 0
