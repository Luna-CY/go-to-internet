package proxy

// ServerConfig 配置结构体
type ServerConfig struct {
    Hostname string // 证书域名
    Port     int    // 监听端口号

    Acme       string // acme工具根路径
    SSLCerFile string // ssl crt路径
    SSLKeyFile string // ssl key路径

    Verbose bool // 详细模式
}
