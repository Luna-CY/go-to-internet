package tunnel

const HandshakeProtocolVersion = 0x02          // 握手协议当前版本
const HandshakeCodeSuccess = 0x00              // 握手响应码
const HandshakeCodeConnectionUpperLimit = 0x01 // 握手响应码
const MessageProtocolVersion = 0x01            // 消息协议当前版本
const CmdNewConnect = 0x01                     // 请求建立目标连接
const CmdData = 0x02                           // 数据传输指令
const MessageCodeNotSet = 0xff                 // 消息响应码
const MessageCodeSuccess = 0x00                // 消息响应码
