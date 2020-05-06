package tunnel

// Config 隧道客户端配置结构体
type Config struct {
    ServerHostname string // 服务端域名
    ServerPort     int    // 服务端端口号
    ServerUsername string // 验证用户
    ServerPassword string // 用户密码

    TargetType     byte   // 目标类型
    TargetHostOrIp string // 目标host或ip
    TargetPort     int    // 目标端口号
}
