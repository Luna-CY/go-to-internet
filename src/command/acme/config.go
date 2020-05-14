package acme

import (
    "flag"
    "fmt"
)

// Config Acme工具配置结构
type Config struct {
    Cmd *flag.FlagSet

    Install bool // 安装acme.sh
}

// Usage 帮助信息
func (c *Config) Usage() {
    _, _ = fmt.Fprintln(c.Cmd.Output(), "manager-go-to-net acme [options]")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "usage:")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net acme -install")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")

    c.Cmd.PrintDefaults()
}

// Validate 检查配置参数是否有效
func (c *Config) Validate() bool {
    return true
}
