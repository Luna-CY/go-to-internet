package user

import (
    "flag"
    "fmt"
)

// Config 用户配置结构
type Config struct {
    Cmd *flag.FlagSet

    List bool // 打印用户列表
    Add  bool // 添加用户
    Upd  bool // 更新用户
    Del  bool // 删除用户

    Config string // 配置文件

    Username string // 用户名
    Password string // 密码
    Expired  string // 有效期，格式: yyyy-MM-dd HH:mm:ss
    MaxRate  int64  // 最大速率，单位kb，格式: 1024 * 1024
}

// Usage 帮助信息
func (c *Config) Usage() {
    _, _ = fmt.Fprintln(c.Cmd.Output(), "manager-go-to-net user [-list|-add|-upd|-del] [options]")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "usage:")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -list")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -add -u test -p password")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -add -u test -p password -e 2020-05-08 23:00:00")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -add -u test -p password -e 2020-05-08 23:00:00 -r 1024 * 1024 * 10")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -upd -u test -p new-password")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net user -del -u test")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")

    c.Cmd.PrintDefaults()
}

// Validate 检查配置参数是否有效
func (c *Config) Validate() bool {
    return true
}
