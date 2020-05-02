package common

// HttpRequest http请求结构体
type HttpRequest struct {
    TargetIp   string `json:"target_ip"`   // 目标ip
    TargetPort int    `json:"target_port"` // 目标端口

    Data []byte `json:"data"` // 数据
}
