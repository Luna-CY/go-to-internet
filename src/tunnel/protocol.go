package tunnel

// 隧道协议
//
// 建立连接
// VER USER_L USER PASS_L PASS
//  1    1     N     1     N
//
// 响应消息
// VER TIME_OUT
//  1     10
//
// 请求建立到目标服务器连接
// VER PORT DST_T DST_N DST
//  1   2     1     1    N
//
// 响应消息
// VER CODE MSG_N MSG
//  1   1     1    N

const VER01 = 0x01
