package main

import (
    "errors"
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/command"
    "gitee.com/Luna-CY/go-to-internet/src/command/acme"
    "gitee.com/Luna-CY/go-to-internet/src/command/user"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "os"
)

func main() {
    // 用户子命令
    userConfig := &user.Config{}
    userCmd := flag.NewFlagSet("user", flag.ExitOnError)
    userConfig.Cmd = userCmd
    userCmd.Usage = userConfig.Usage

    userCmd.BoolVar(&userConfig.List, "list", false, "打印用户列表")
    userCmd.BoolVar(&userConfig.Add, "add", false, "新增用户，该动作需要提供用户名以及用户密码")
    userCmd.BoolVar(&userConfig.Upd, "upd", false, "更新用户，该动作需要提供用户名以及任意其他配置项")
    userCmd.BoolVar(&userConfig.Del, "del", false, "删除用户，该动作需要提供用户名")

    userCmd.StringVar(&userConfig.Config, "c", "/etc/go-to-net/users.json", "配置文件位置")
    userCmd.StringVar(&userConfig.Username, "u", "", "用户名，用户的唯一标识")
    userCmd.StringVar(&userConfig.Password, "p", "", "用户密码")
    userCmd.StringVar(&userConfig.Expired, "expired", "", "用户过期时间，-表示不过期，格式: 2006-01-02T15:04:05 (default -)")
    userCmd.IntVar(&userConfig.MaxRate, "max-rate", -1, "用户最大传输速率，-1表示未设置该参数；0表示不限速，非0表示最大传输速率，单位为KB")
    userCmd.IntVar(&userConfig.MaxConnection, "max-connection", -1, "用户最大连接数，-1表示未设置该参数；0表示不限制连接数，非0表示最大连接数")

    acmeConfig := &acme.Config{}
    acmeCmd := flag.NewFlagSet("acme", flag.ExitOnError)
    acmeConfig.Cmd = acmeCmd
    acmeCmd.Usage = acmeConfig.Usage

    acmeCmd.BoolVar(&acmeConfig.Install, "install", false, "安装acme工具")
    acmeCmd.BoolVar(&acmeConfig.Issue, "issue", false, "申请证书")

    acmeCmd.StringVar(&acmeConfig.Hostname, "hostname", "", "操作的域名")

    if len(os.Args) < 2 || "-h" == os.Args[1] || "--help" == os.Args[1] {
        _, _ = fmt.Fprintf(flag.CommandLine.Output(), "version %v\n", common.Version)
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "manager-go-to-net subcommand options")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "sub commands:")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    user: 用户管理命令")
        _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    acme: Acme证书工具管理命令")

        os.Exit(0)
    }

    switch os.Args[1] {
    case "user":
        if err := parse(userCmd, userConfig); nil != err {
            userCmd.Usage()

            return
        }

        if err := user.Exec(userConfig); nil != err {
            logger.Errorf("处理操作失败: %v", err)
        }
    case "acme":
        if err := parse(acmeCmd, acmeConfig); nil != err {
            acmeCmd.Usage()

            return
        }

        if err := acme.Exec(acmeConfig); nil != err {
            logger.Errorf("处理操作失败: %v", err)
        }
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
