package main

import (
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/socket"
    "os"
)

// serverCommandUsage 打印控制台Usage信息
func clientCommandUsage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "cli-go-to-net -sh Hostname -u USERNAME -p PASSWORD [options]")

    flag.PrintDefaults()
}

func main() {
    server := socket.Socket{}

    flag.StringVar(&server.Hostname, "sh", "", "服务器域名")
    flag.IntVar(&server.Port, "sp", 443, "服务器端口")
    flag.StringVar(&server.LocalAddr, "la", "127.0.0.1", "本地监听地址")
    flag.IntVar(&server.LocalPort, "lp", 1280, "本地监听端口")

    flag.StringVar(&server.Username, "u", "", "服务端身份认证用户名")
    flag.StringVar(&server.Password, "p", "", "服务端身份认证密码")

    flag.BoolVar(&server.Verbose, "v", false, "打印详细日志")

    flag.Usage = clientCommandUsage
    flag.Parse()

    if "" == server.Hostname || "" == server.Username || "" == server.Password {
        flag.Usage()

        os.Exit(0)
    }

    server.Start()
}
