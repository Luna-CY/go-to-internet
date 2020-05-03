package common

// HttpRequest http请求结构体
type HttpRequest struct {
    TargetIp   string `json:"target_ip"`   // 目标ip
    TargetPort int    `json:"target_port"` // 目标端口

    Username string `json:"username"` // 鉴权用户
    Password string `json:"password"` // 鉴权密码

    Data []byte `json:"data"` // 数据
}

// HttpResponse http响应结构体
type HttpResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
