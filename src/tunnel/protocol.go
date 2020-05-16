package tunnel

// 隧道协议
//
// 建立连接
// VER USER_L USER PASS_L PASS PORT DST_T DST_N DST
//  1    1     N     1     N    2     1     1    N
//
// 响应消息
// VER CODE MSG_N MSG
//  1   1     1    N

const Ver01 = 0x01

const Success = 0x00
const SuccessMessage = "OK"

const ConnectionUpperLimit = 0x01
const ConnectionUpperLimitMessage = "已到达连接上限"

const CheckConnectTargetIp = "0.0.0.0"
const CheckConnectTargetPort = 0
