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
        return utils.ExecCommandOutputToLog("systemctl", []string{"start", "go-to-net"}, nil)
    case c.Config.Stop:
        return utils.ExecCommandOutputToLog("systemctl", []string{"stop", "go-to-net"}, nil)
    case c.Config.Enable:
        return utils.ExecCommandOutputToLog("systemctl", []string{"enable", "go-to-net"}, nil)
    case c.Config.Disable:
        return utils.ExecCommandOutputToLog("systemctl", []string{"disable", "go-to-net"}, nil)
    default:
        return errors.New("无效命令")
    }
}

// generateServiceConfig 生成服务配置文件
func (c *Cmd) generateServiceConfig() error {
    if "" == c.Config.Hostname {
        return errors.New("域名不能为空")
    }

    exec := c.Config.ExecCmd
    if "" == exec {
        result, err := utils.ExecCommandGetStdout("which", []string{"ser-go-to-net"}, nil)
        if nil != err || 0 == len(result) {
            return errors.New("未找到ser-go-to-net命令，可以使用 -exec 参数指定ser-go-to-net命令位置")
        }
        exec = result[0]
    }

    filepath := "/etc/systemd/system/go-to-net.service"
    file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if nil != err {
        return err
    }

    content := strings.Replace(template, "EXEC_CMD", exec, 1)
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

    return utils.ExecCommandOutputToLog("systemctl", []string{"daemon-reload"}, nil)
}
