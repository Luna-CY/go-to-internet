package main

import (
    "crypto/tls"
    "errors"
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/proxy"
    "os"
    "os/user"
    "path"
    "strconv"
    "syscall"
)

// serverCommandUsage 打印控制台Usage信息
func serverCommandUsage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "ser-go-to-net -H Hostname [options]")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "此工具默认通过acme.sh工具来管理证书，如果不使用acme.sh工具，需要设置-c与-k参数指定证书位置")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "usage:")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    ser-go-to-net -H proxy.example.com")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    ser-go-to-net -H proxy.example.com -acme /root/.acme.sh")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "    ser-go-to-net -H proxy.example.com -c /path/to/cert.cer -k /path/to/key.key")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "")

    flag.PrintDefaults()
}

func main() {
    config := &proxy.ServerConfig{}

    flag.StringVar(&config.Hostname, "H", "", "域名，该域名应该与证书的域名一致")
    flag.IntVar(&config.Port, "p", 443, "监听端口号")

    flag.StringVar(&config.Acme, "acme", "/root/.acme.sh", "acme工具的根路径")
    flag.StringVar(&config.SSLCerFile, "c", "", "SSL CER文件路径 (default \"${arg:acme}/${arg:H}/fullchain.cer\")")
    flag.StringVar(&config.SSLKeyFile, "k", "", "SSL KEY文件路径 (default \"${arg:acme}/${arg:H}/${arg:H}.key\")")

    flag.StringVar(&config.User, "u", "", "设置运行时用户")
    flag.StringVar(&config.Group, "g", "", "设置运行时用户组")

    flag.BoolVar(&config.Verbose, "v", false, "打印详细日志")

    flag.Usage = serverCommandUsage
    flag.Parse()

    if "" == config.Hostname || ("" == config.Acme && ("" == config.SSLCerFile || "" == config.SSLKeyFile)) {
        flag.Usage()

        os.Exit(0)
    }

    if "" != config.User {
        if err := switchToUser(config.User); nil != err {
            logger.Errorf("切换用户失败: %v", err)

            os.Exit(1)
        }
    }

    if "" != config.Group {
        if err := switchToGroup(config.Group); nil != err {
            logger.Errorf("切换用户组失败: %v", err)

            os.Exit(1)
        }
    }

    tlsListen(config)
}

// switchToUser 切换运行时用户
func switchToUser(username string) error {
    info, err := user.Lookup(username)
    if nil != err {
        return err
    }

    uid, err := strconv.Atoi(info.Uid)
    if nil != err {
        return err
    }

    if err := syscall.Setuid(uid); nil != err {
        return err
    }

    gid, err := strconv.Atoi(info.Gid)
    if nil != err {
        return err
    }

    return syscall.Setgid(gid)
}

// switchToGroup 切换运行时用户组
func switchToGroup(group string) error {
    info, err := user.LookupGroup(group)
    if nil != err {
        return err
    }

    gid, err := strconv.Atoi(info.Gid)
    if nil != err {
        return err
    }

    return syscall.Setgid(gid)
}

// tlsListen 启动tls服务器
func tlsListen(config *proxy.ServerConfig) {
    tlsConfig, err := getTlsConfig(config)
    if nil != err {
        logger.Errorf("加载TLS证书失败: %v", err)

        return
    }

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

        go proxy.StartConnection(conn, config.Hostname, config.Verbose)
    }
}

// getTlsConfig 获取TSL配置结构
func getTlsConfig(config *proxy.ServerConfig) (*tls.Config, error) {
    cert, key := config.SSLCerFile, config.SSLKeyFile

    if "" == cert || "" == key {
        dir, err := os.Stat(config.Acme)
        if nil != err {
            return nil, errors.New(fmt.Sprintf("%v: 无法找到该路径", config.Acme))
        }

        if !dir.IsDir() {
            return nil, errors.New("acme路径不是一个目录")
        }

        cert = path.Join(config.Acme, config.Hostname, "fullchain.cer")
        key = path.Join(config.Acme, config.Hostname, fmt.Sprintf("%v.key", config.Hostname))
    }

    certificate, err := tls.LoadX509KeyPair(cert, key)
    if nil != err {
        return nil, err
    }

    return &tls.Config{Certificates: []tls.Certificate{certificate}}, nil
}
