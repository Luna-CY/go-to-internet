package config

// Server 配置对象
type Server struct {
    Hostname string // 证书域名
    Port     int    // 监听端口号

    SSLCrtFile string // ssl crt路径
    SSLKeyFile string // ssl key路径

    Authorize bool // 是否验证用户身份
}
