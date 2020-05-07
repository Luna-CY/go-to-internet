package main

import (
    "crypto/tls"
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/proxy"
    "os"
)

// serverCommandUsage 打印控制台Usage信息
func serverCommandUsage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "server -H Hostname -c CRT -k KEY [options]")

    flag.PrintDefaults()
}

func main() {
    config := &proxy.ServerConfig{}

    flag.StringVar(&config.Hostname, "H", "", "域名，该域名应该与证书的域名一致")
    flag.IntVar(&config.Port, "p", 443, "监听端口号")

    flag.StringVar(&config.SSLCrtFile, "c", "", "SSL CRT文件路径")
    flag.StringVar(&config.SSLKeyFile, "k", "", "SSL KEY文件路径")

    flag.BoolVar(&config.Verbose, "v", false, "打印详细日志")

    flag.Usage = serverCommandUsage
    flag.Parse()

    if "" == config.Hostname || "" == config.SSLCrtFile || "" == config.SSLKeyFile {
        flag.Usage()

        os.Exit(0)
    }

    tlsListen(config)
}

// tlsListen 启动tls服务器
func tlsListen(config *proxy.ServerConfig) {
    cert, err := tls.LoadX509KeyPair(config.SSLCrtFile, config.SSLKeyFile)
    if nil != err {
        logger.Errorf("加载TLS证书失败: %v", err)

        return
    }
    tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
    ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", config.Port), tlsConfig)
    if nil != err {
        logger.Errorf("启动服务器失败: %v", err)

        return
    }
    defer ln.Close()
    logger.Infof("启动监听 %v:%d ...\n", config.Hostname, config.Port)

    for {
        conn, err := ln.Accept()
        if nil != err {
            continue
        }

        go proxy.StartConnection(conn, config.Verbose)
    }
}
