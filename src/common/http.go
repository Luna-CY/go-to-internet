package common

// HttpRequest http请求结构体
type HttpRequest struct {
    TargetIp   string `json:"target_ip"`   // 目标ip
    TargetPort int    `json:"target_port"` // 目标端口

    Data []byte `json:"data"` // 数据
}

// HttpResponse http响应结构体
type HttpResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`

    Data   []byte `json:"data,omitempty"`
    IsLast bool   `json:"is_last,omitempty"`
}
