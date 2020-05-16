package acme

import (
    "flag"
    "fmt"
)

// Config Acme工具配置结构
type Config struct {
    Cmd *flag.FlagSet

    Install bool // 安装acme.sh

    Issue      bool // 申请证书
    Standalone bool // 模拟一个http服务器
    Nginx      bool // 通过nginx验证

    Hostname string // 操作的域名
}

// Usage 帮助信息
func (c *Config) Usage() {
    _, _ = fmt.Fprintln(c.Cmd.Output(), "manager-go-to-net acme [options]")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "usage:")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net acme -install")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net acme -issue -hostname YOUR_HOST")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net acme -renew -hostname YOUR_HOST")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")

    c.Cmd.PrintDefaults()
}

// Validate 检查配置参数是否有效
func (c *Config) Validate() bool {
    switch {
    case c.Install:
        return true
    case c.Issue:
        if "" == c.Hostname {
            return false
        }

        return true
    default:
        return false
    }
}
