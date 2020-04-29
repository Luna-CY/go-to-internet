package main

import (
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/proxy"
    "os"
)

// serverCommandUsage 打印控制台Usage信息
func clientCommandUsage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "client -H Hostname [-auth -u USERNAME -P PASSWORD] [options]")

    flag.PrintDefaults()
}

func main() {
    socket := proxy.Socket{}

    flag.StringVar(&socket.Hostname, "H", "", "服务器域名")
    flag.StringVar(&socket.IpAddr, "l", "127.0.0.1", "本地监听地址")
    flag.IntVar(&socket.Port, "p", 1280, "本地监听端口")

    flag.BoolVar(&socket.Authorize, "auth", false, "服务端是否需要身份认证")
    flag.StringVar(&socket.Username, "u", "", "服务端身份认证用户名")
    flag.StringVar(&socket.Password, "P", "", "服务端身份认证密码")

    flag.Usage = clientCommandUsage
    flag.Parse()

    if "" == socket.Hostname || (socket.Authorize && ("" == socket.Username || "" == socket.Password)) {
        flag.Usage()

        os.Exit(0)
    }

    socket.Start()
}
