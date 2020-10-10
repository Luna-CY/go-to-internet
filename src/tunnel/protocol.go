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
// VER CMD CODE DATA_L DATA
//  1   1   1     2     N
//
// 新建连接的DATA部分协议
//
// IP_TYPE PORT IP
//    1     2   N
