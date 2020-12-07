package service

import (
    "flag"
    "fmt"
)

// Config service命令的配置结构体
type Config struct {
    Cmd *flag.FlagSet

    Install        bool   // 安装到系统服务
    Start          bool   // 启动服务
    Stop           bool   // 关闭服务
    Enable         bool   // 设置开机自启动
    Disable        bool   // 取消开机自启动
    Client         bool   // 安装客户端服务
    Username       string // 客户端服务的用户名称
    Password       string // 客户端服务的用户密码
    Restart        bool   // 重启服务
    SetAutoRestart bool   // 自动重启服务
    Remove         bool   // 从系统服务移除

    Hostname string // 域名
    ExecCmd  string // ser-go-to-net可执行文件位置
    Cron     string // 计划任务周期
}

func (c *Config) Usage() {
    _, _ = fmt.Fprintln(c.Cmd.Output(), "manager-go-to-net service [-install|-start|-stop|-restart|-enable|-disable|-auto-restart|-remove] [options]")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "该命令支持通过默认参数安装代理服务到系统")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "usage:")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -install -hostname YOUR_HOST -exec EXEC_PATH")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -start")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -stop")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -set-auto-restart -cron '0 6 * * *'")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "    manager-go-to-net service -remove")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")
    _, _ = fmt.Fprintln(c.Cmd.Output(), "")

    c.Cmd.PrintDefaults()
}

func (c *Config) Validate() bool {
    if !c.Install && !c.Start && !c.Stop && !c.Restart && !c.Enable && !c.Disable && !c.SetAutoRestart && !c.Remove {
        return false
    }

    if c.Install && ("" == c.Hostname || "" == c.ExecCmd) {
        return false
    }

    if c.Client && ("" == c.Username || "" == c.Password) {
        return false
    }

    return true
}
