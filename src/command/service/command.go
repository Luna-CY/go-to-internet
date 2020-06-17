package service

import (
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "os"
    "strings"
)

type Cmd struct {
    Config *Config
}

func (c *Cmd) Exec() error {
    switch {
    case c.Config.Install:
        return c.generateServiceConfig()
    case c.Config.Start:
    case c.Config.Stop:
    case c.Config.Enable:
    case c.Config.Disable:
    default:
        return errors.New("无效命令")
    }

    return nil
}

// generateServiceConfig 生成服务配置文件
func (c *Cmd) generateServiceConfig() error {
    if "" == c.Config.Hostname {
        return errors.New("域名不能为空")
    }

    result, err := utils.ExecCommandGetStdout("which", []string{"ser-go-to-net"}, nil)
    if nil != err {
        return err
    }

    if 0 == len(result) {
        return errors.New("未找到ser-go-to-net命令")
    }

    filepath := "/etc/systemd/system/go-to-net.service"
    file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if nil != err {
        return err
    }

    content := strings.Replace(template, "EXEC_CMD", result[0], 1)
    content = strings.Replace(content, "YOUR_HOST", c.Config.Hostname, 1)

    n, err := file.Write([]byte(content))
    if nil != err {
        return err
    }

    if n != len(content) {
        if err := os.Remove(filepath); nil != err {
            return errors.New(fmt.Sprintf("生成服务配置文件失败: %v", err))
        }

        return errors.New("生成服务配置文件失败")
    }

    return nil
}
