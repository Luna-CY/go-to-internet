package main

import (
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/socket"
    "os"
)

// serverCommandUsage 打印控制台Usage信息
func clientCommandUsage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "client -H Hostname -u USERNAME -P PASSWORD [options]")

    flag.PrintDefaults()
}

func main() {
    server := socket.Socket{}

    flag.StringVar(&server.Hostname, "H", "", "服务器域名")
    flag.IntVar(&server.Port, "p", 443, "服务器端口")
    flag.StringVar(&server.LocalAddr, "l", "127.0.0.1", "本地监听地址")
    flag.IntVar(&server.LocalPort, "lp", 1280, "本地监听端口")

    flag.StringVar(&server.Username, "u", "", "服务端身份认证用户名")
    flag.StringVar(&server.Password, "P", "", "服务端身份认证密码")

    flag.BoolVar(&server.Verbose, "v", false, "打印详细日志")

    flag.Usage = clientCommandUsage
    flag.Parse()

    if "" == server.Hostname || "" == server.Username || "" == server.Password {
        flag.Usage()

        os.Exit(0)
    }

    server.Start()
}
