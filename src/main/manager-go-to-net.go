package main

import (
    "flag"
    "gitee.com/Luna-CY/go-to-internet/src/command/user"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "os"
)

func main() {
    userConfig := &user.Config{}

    // 用户子命令
    userCmd := flag.NewFlagSet("user", flag.ExitOnError)
    userCmd.BoolVar(&userConfig.List, "list", false, "打印用户列表")
    userCmd.BoolVar(&userConfig.Add, "add", false, "新增用户，该动作需要提供所有用户配置项")
    userCmd.BoolVar(&userConfig.Upd, "upd", false, "更新用户，该动作需要提供用户名以及任意其他配置项")
    userCmd.BoolVar(&userConfig.Del, "del", false, "删除用户，该动作需要提供用户名")

    userCmd.StringVar(&userConfig.Username, "u", "", "用户名，用户的唯一标识")
    userCmd.StringVar(&userConfig.Password, "p", "", "用户密码")
    userCmd.StringVar(&userConfig.Expired, "e", "-", "用户过期时间，单短横线代表不过期")
    userCmd.Int64Var(&userConfig.MaxRate, "r", 0, "用户最大速率，0代表不限速，单位为kb")

    if len(os.Args) < 2 {
        logger.Error("必须提供子命令，可用子命令请查看帮助信息")

        return
    }

    switch os.Args[1] {
    case "user":
        if err := userCmd.Parse(os.Args[2:]); nil != err {
            logger.Errorf("解析命令失败: %v", err)

            return
        }

        logger.Info(userConfig)
    default:
        logger.Error("不支持的子命令，请查看帮助信息")

        return
    }
}
