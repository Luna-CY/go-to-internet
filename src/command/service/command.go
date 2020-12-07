package service

import (
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "io/ioutil"
    "os"
    "strings"
)

type Cmd struct {
    Config Config
}

func (c *Cmd) Exec() error {
    switch {
    case c.Config.Install:
        return c.installService()
    case c.Config.Start:
        return utils.ExecCommandOutputToLog("systemctl", []string{"start", "go-to-net"}, nil)
    case c.Config.Stop:
        return utils.ExecCommandOutputToLog("systemctl", []string{"stop", "go-to-net"}, nil)
    case c.Config.Restart:
        return utils.ExecCommandOutputToLog("systemctl", []string{"restart", "go-to-net"}, nil)
    case c.Config.Enable:
        return utils.ExecCommandOutputToLog("systemctl", []string{"enable", "go-to-net"}, nil)
    case c.Config.Disable:
        return utils.ExecCommandOutputToLog("systemctl", []string{"disable", "go-to-net"}, nil)
    case c.Config.SetAutoRestart:
        if "" == c.Config.Cron {
            c.Config.Cron = "0 6 1 * *"
        }

        return c.setCron(c.Config.Cron)
    case c.Config.Remove:
        return c.removeService()
    default:
        return errors.New("无效命令")
    }
}

// installService 安装服务
func (c *Cmd) installService() error {
    if "" == c.Config.Hostname {
        return errors.New("域名不能为空")
    }

    filepath := "/etc/systemd/system/go-to-net.service"
    file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if nil != err {
        return err
    }
    defer file.Close()

    var template string
    if c.Config.Client {
        template = clientTemplate
    } else {
        template = serverTemplate
    }

    content := strings.Replace(template, "EXEC_CMD", c.Config.ExecCmd, 1)
    content = strings.Replace(content, "YOUR_HOST", c.Config.Hostname, 1)
    content = strings.Replace(content, "USERNAME", c.Config.Username, 1)
    content = strings.Replace(content, "PASSWORD", c.Config.Password, 1)

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

    if err := utils.ExecCommandOutputToLog("systemctl", []string{"daemon-reload"}, nil); nil != err {
        return err
    }

    if "-" != c.Config.Cron {
        if "" == c.Config.Cron {
            c.Config.Cron = "0 6 1 * *"
        }

        if err := c.setCron(c.Config.Cron); nil != err {
            return err
        }
    }

    return nil
}

// setCron 设置计划任务
func (c *Cmd) setCron(cron string) error {
    filepath := "/etc/crontab"
    if _, err := os.Stat(filepath); nil != err {
        return err
    }

    file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
    if nil != err {
        return err
    }
    defer file.Close()

    dataBytes, err := ioutil.ReadAll(file)
    if nil != err {
        return err
    }

    var result []string

    content := string(dataBytes)
    rows := strings.Split(content, "\n")
    for _, row := range rows {
        if "" == row {
            continue
        }

        // min hour day month week user cmd
        columns := strings.Split(row, " ")

        // 长度小于7是错误配置，跳过不管
        if 7 > len(columns) {
            result = append(result, strings.Join(columns, " "))

            continue
        }

        // 跳过这个任务
        if "systemctl restart go-to-net" == strings.Join(columns[6:], " ") {
            continue
        }

        result = append(result, strings.Join(columns, " "))
    }

    if err := file.Truncate(0); nil != err {
        return err
    }

    if _, err := file.Seek(0, 0); nil != err {
        return err
    }

    // 添加计划任务
    if "-" != cron {
        result = append(result, fmt.Sprintf("%v root systemctl restart go-to-net", c.Config.Cron))
    }

    writeContent := []byte(strings.TrimSpace(strings.Join(result, "\n")))
    if l, err := file.Write(writeContent); nil != err || l != len(writeContent) {
        return errors.New(fmt.Sprintf("写入文件失败，总长度: %v 写入长度: %v 错误信息: %v", len(writeContent), l, err))
    }

    return nil
}

// removeService 移除移动服务
func (c *Cmd) removeService() error {
    if err := utils.ExecCommandOutputToLog("systemctl", []string{"stop", "go-to-net"}, nil); nil != err {
        return err
    }

    if err := utils.ExecCommandOutputToLog("systemctl", []string{"disable", "go-to-net"}, nil); nil != err {
        return err
    }

    filepath := "/etc/systemd/system/go-to-net.service"
    _, err := os.Stat(filepath)
    if nil != err {
        if os.IsNotExist(err) {
            return nil
        }

        return err
    }

    if err := os.Remove(filepath); nil != err {
        return err
    }

    return c.setCron("-")
}
