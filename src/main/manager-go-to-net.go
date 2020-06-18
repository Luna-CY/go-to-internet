package main

import (
    "errors"
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/command"
    "gitee.com/Luna-CY/go-to-internet/src/command/acme"
    "gitee.com/Luna-CY/go-to-internet/src/command/service"
    "gitee.com/Luna-CY/go-to-internet/src/command/user"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "os"
)

func main() {
    if len(os.Args) < 2 || "-h" == os.Args[1] || "--help" == os.Args[1] {
        _, _ = fmt.Fprintf(flag.CommandLine.Output(), "version %v\n", common.Version)
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "manager-go-to-net subcommand options")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "sub commands:")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    user: 用户管理命令")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    acme: Acme证书工具管理命令")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    service: 服务管理命令")

        os.Exit(0)
    }

    switch os.Args[1] {
    case "user":
        userSubCommand()
    case "acme":
        acmeSubCommand()
    case "service":
        serviceSubCommand()
    default:
        logger.Error("不支持的子命令，请查看帮助信息")

        return
    }
}

// parse 解析命令行输入
func parse(cmd *flag.FlagSet, config command.Config) error {
    if err := cmd.Parse(os.Args[2:]); nil != err {
        return errors.New(fmt.Sprintf("解析命令失败: %v", err))
    }

    if !config.Validate() {
        return errors.New("校验参数失败")
    }

    return nil
}

// userSubCommand 用户子命令
func userSubCommand() {
    config := &user.Config{}
    userCmd := flag.NewFlagSet("user", flag.ExitOnError)
    config.Cmd = userCmd
    userCmd.Usage = config.Usage

    userCmd.BoolVar(&config.List, "list", false, "打印用户列表")
    userCmd.BoolVar(&config.Add, "add", false, "新增用户，该动作需要提供用户名以及用户密码")
    userCmd.BoolVar(&config.Upd, "upd", false, "更新用户，该动作需要提供用户名以及任意其他配置项")
    userCmd.BoolVar(&config.Del, "del", false, "删除用户，该动作需要提供用户名")

    userCmd.StringVar(&config.Config, "c", "/etc/go-to-net/users.json", "配置文件位置")
    userCmd.StringVar(&config.Username, "u", "", "用户名，用户的唯一标识")
    userCmd.StringVar(&config.Password, "p", "", "用户密码")
    userCmd.StringVar(&config.Expired, "expired", "", "用户过期时间，-表示不过期，格式: 2006-01-02T15:04:05 (default -)")
    userCmd.IntVar(&config.MaxRate, "max-rate", -1, "用户最大传输速率，-1表示未设置该参数；0表示不限速，非0表示最大传输速率，单位为KB")
    userCmd.IntVar(&config.MaxConnection, "max-connection", -1, "用户最大连接数，-1表示未设置该参数；0表示不限制连接数，非0表示最大连接数")

    if err := parse(userCmd, config); nil != err {
        userCmd.Usage()

        return
    }

    cmd := &user.Cmd{Config: config}
    if err := cmd.Exec(); nil != err {
        logger.Errorf("处理操作失败: %v", err)
    }
}

// acmeSubCommand acme子命令
func acmeSubCommand() {
    config := &acme.Config{}
    acmeCmd := flag.NewFlagSet("acme", flag.ExitOnError)
    config.Cmd = acmeCmd
    acmeCmd.Usage = config.Usage

    acmeCmd.BoolVar(&config.Install, "install", false, "安装acme工具")
    acmeCmd.BoolVar(&config.Issue, "issue", false, "申请证书")
    acmeCmd.BoolVar(&config.Nginx, "nginx", false, "通过nginx验证服务器")
    acmeCmd.BoolVar(&config.Standalone, "standalone", false, "通过acme.sh的standalone模式验证服务器")

    acmeCmd.StringVar(&config.Hostname, "hostname", "", "操作的域名")

    if err := parse(acmeCmd, config); nil != err {
        acmeCmd.Usage()

        return
    }

    cmd := &acme.Cmd{Config: config}
    if err := cmd.Exec(); nil != err {
        logger.Errorf("处理操作失败: %v", err)
    }
}

// serviceSubCommand 服务子命令
func serviceSubCommand() {
    config := &service.Config{}
    serviceCmd := flag.NewFlagSet("service", flag.ExitOnError)
    config.Cmd = serviceCmd
    serviceCmd.Usage = config.Usage

    serviceCmd.BoolVar(&config.Install, "install", false, "安装服务到系统")
    serviceCmd.BoolVar(&config.Start, "start", false, "启动服务")
    serviceCmd.BoolVar(&config.Stop, "stop", false, "关闭服务")
    serviceCmd.BoolVar(&config.Enable, "enable", false, "添加到开机启动")
    serviceCmd.BoolVar(&config.Disable, "disable", false, "取消开机启动")

    serviceCmd.StringVar(&config.Hostname, "hostname", "", "域名，指定-install项时必须指定此项")
    serviceCmd.StringVar(&config.ExecCmd, "exec", "", "ser-go-to-net命令位置，指定-install项时必须指定此项")

    if err := parse(serviceCmd, config); nil != err {
        serviceCmd.Usage()

        return
    }

    cmd := &service.Cmd{Config: config}
    if err := cmd.Exec(); nil != err {
        logger.Errorf("处理操作失败: %v", err)
    }
}
