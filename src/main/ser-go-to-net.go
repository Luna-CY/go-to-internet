package main

import (
    "crypto/tls"
    "encoding/json"
    "errors"
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/config"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/proxy"
    "golang.org/x/sys/unix"
    "io/ioutil"
    "os"
    "os/user"
    "path"
    "strconv"
)

// serverCommandUsage 打印控制台Usage信息
func serverCommandUsage() {
    _, _ = fmt.Fprintf(flag.CommandLine.Output(), "version %v\n", common.Version)
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
    serverConfig := &proxy.Config{}

    flag.StringVar(&serverConfig.Hostname, "H", "", "域名，该域名应该与证书的域名一致")
    flag.IntVar(&serverConfig.Port, "p", 443, "监听端口号")

    flag.StringVar(&serverConfig.Acme, "acme", "/root/.acme.sh", "acme工具的根路径")
    flag.StringVar(&serverConfig.SSLCerFile, "c", "", "SSL CER文件路径 (default \"${arg:acme}/${arg:H}/fullchain.cer\")")
    flag.StringVar(&serverConfig.SSLKeyFile, "k", "", "SSL KEY文件路径 (default \"${arg:acme}/${arg:H}/${arg:H}.key\")")

    // golang 暂时支持不了切换运行时用户，fork子进程的成本又比较高不值得
    //flag.StringVar(&serverConfig.User, "u", "", "设置运行时用户")
    //flag.StringVar(&serverConfig.Group, "g", "", "设置运行时用户组")

    flag.StringVar(&serverConfig.UserConfig, "uc", "/etc/go-to-net/users.json", "用户配置文件，可以通过manager-go-to-net命令生成")

    flag.BoolVar(&serverConfig.Verbose, "v", false, "打印详细日志")

    flag.Usage = serverCommandUsage
    flag.Parse()

    if "" == serverConfig.Hostname || "" == serverConfig.UserConfig || ("" == serverConfig.Acme && ("" == serverConfig.SSLCerFile || "" == serverConfig.SSLKeyFile)) {
        flag.Usage()

        os.Exit(0)
    }

    if "" != serverConfig.User {
        if err := switchToUser(serverConfig.User); nil != err {
            logger.Errorf("切换用户失败: %v", err)

            os.Exit(1)
        }
    }

    if "" != serverConfig.Group {
        if err := switchToGroup(serverConfig.Group); nil != err {
            logger.Errorf("切换用户组失败: %v", err)

            os.Exit(1)
        }
    }

    tlsListen(serverConfig)
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

    if err := unix.Setuid(uid); nil != err {
        return err
    }

    gid, err := strconv.Atoi(info.Gid)
    if nil != err {
        return err
    }

    return unix.Setgid(gid)
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

    return unix.Setgid(gid)
}

// tlsListen 启动tls服务器
func tlsListen(config *proxy.Config) {
    userConfig, err := loadUserConfig(config.UserConfig)
    if nil != err {
        logger.Error(err)

        return
    }

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
    logger.Infof("启动监听 %v:%d ...", config.Hostname, config.Port)

    proxyInstance := proxy.Proxy{UserConfig: userConfig, Hostname: config.Hostname, Verbose: config.Verbose}
    if err := proxyInstance.Init(); nil != err {
        logger.Errorf("初始化服务器失败: %v", err)

        return
    }

    for {
        conn, err := ln.Accept()
        if nil != err {
            continue
        }

        go proxyInstance.Accept(conn)
    }
}

// getTlsConfig 获取TSL配置结构
func getTlsConfig(config *proxy.Config) (*tls.Config, error) {
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

// loadUserConfig 加载用户配置
func loadUserConfig(c string) (*config.UserConfig, error) {
    file, err := os.Open(c)
    if nil != err {
        return nil, errors.New(fmt.Sprintf("无法打开配置文件: %v", err))
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if nil != err {
        return nil, errors.New(fmt.Sprintf("无法读取配置文件: %v", err))
    }

    userConfig := &config.UserConfig{}
    if err := json.Unmarshal(data, userConfig); nil != err {
        return nil, errors.New(fmt.Sprintf("解析配置文件失败: %v", err))
    }

    return userConfig, nil
}
