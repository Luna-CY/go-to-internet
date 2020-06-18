package service

import "C"
import (
    "flag"
    "fmt"
)

// Config service命令的配置结构体
type Config struct {
    Cmd *flag.FlagSet

    Install bool // 安装到系统服务
    Start   bool // 启动服务
    Stop    bool // 关闭服务
    Enable  bool // 设置开机自启动
    Disable bool // 取消开机自启动

    Hostname string // 域名
    ExecCmd  string // ser-go-to-net可执行文件位置
}

func (c *Config) Usage() {
    _, _ = fmt.Fprintln(c.Cmd.Output(), "manager-go-to-net service [-install|-start|-stop|-enable|-disable] [options]")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "usage:")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -install -hostname YOUR_HOST")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -start")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -stop")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")

    c.Cmd.PrintDefaults()
}

func (c *Config) Validate() bool {
    if !c.Install && !c.Start && !c.Stop && !c.Enable && !c.Disable {
        return false
    }

    if c.Install && "" == c.Hostname {
        return false
    }

    return true
}
