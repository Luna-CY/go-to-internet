package acme

import (
    "bufio"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/common"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "gitee.com/Luna-CY/go-to-internet/src/utils"
    "golang.org/x/sys/unix"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path"
    "strings"
)

type Cmd struct {
    Config *Config
}

func (c *Cmd) Exec() error {
    switch {
    case c.Config.Install:
        home, err := os.UserHomeDir()
        if nil != err {
            logger.Error("无法获取用户主目录")

            break
        }
        acmePath := path.Join(home, ".acme.sh")

        info, err := os.Stat(acmePath)
        if nil == err {
            if !info.IsDir() {
                logger.Errorf("[%v]路径已存在并且不是一个目录", acmePath)

                break
            }

            logger.Infof("已安装acme.sh工具，重新安装请删除[%v]目录后重新执行", acmePath)

            break
        }

        if !os.IsNotExist(err) {
            logger.Errorf("获取路径信息失败: %v", err)

            break
        }

        output := path.Join(os.TempDir(), "install-acme.sh")

        if err := c.download(common.AcmePath, output); nil != err {
            logger.Errorf("下载安装脚本失败: %v", err)

            break
        }

        if err := os.Chmod(output, os.FileMode(0755)); nil != err {
            logger.Errorf("修改文件权限失败: %v", err)

            break
        }

        if err := utils.ExecCommandOutputToLog("sh", []string{"-c", output}, &[]string{"INSTALLONLINE=1"}); nil != err {
            logger.Errorf("安装acme.sh失败: %v", err)

            break
        }

        if err := unix.Unlink(output); nil != err {
            logger.Errorf("删除安装脚本失败: %v", err)

            break
        }

        logger.Info("安装完成")
    case c.Config.Issue:
        home, err := os.UserHomeDir()
        if nil != err {
            logger.Error("无法获取用户主目录")

            break
        }

        command := path.Join(home, ".acme.sh", "acme.sh")
        info, err := os.Stat(command)
        if nil != err {
            logger.Errorf("没有找到acme.sh工具: %v", err)

            break
        }

        if info.IsDir() {
            logger.Error("无效的acme.sh工具路径")

            break
        }

        switch {
        case c.Config.Nginx:
            if err := c.checkAndInstallNginx(); nil != err {
                logger.Error(err)

                break
            }

            if err := c.generateNginxConfig(c.Config.Hostname); nil != err {
                logger.Errorf("创建nginx配置失败: %v", err)

                break
            }

            isExist, err := utils.FileExists(fmt.Sprintf("/root/.acme.sh/%v", c.Config.Hostname))
            if nil != err {
                logger.Errorf("检查证书路径失败: %v", err)

                break
            }

            if isExist {
                logger.Info("该域名已存在证书")

                break
            }

            if err := utils.ExecCommandOutputToLog(command, []string{"--issue", "-d", c.Config.Hostname, "--nginx"}, nil); nil != err {
                logger.Errorf("申请证书失败: %v", err)

                break
            }

            logger.Info("申请证书完成")
        default:
            if err := utils.ExecCommandOutputToLog(command, []string{"--issue", "-d", c.Config.Hostname, "--standalone"}, nil); nil != err {
                logger.Errorf("申请证书失败: %v", err)

                break
            }

            logger.Info("申请证书完成")
        }
    }

    return nil
}

// checkAndInstallNginx 检查nginx是否存在，不存在时安装nginx
func (c *Cmd) checkAndInstallNginx() error {
    cmd := exec.Command("nginx", "-v")
    if err := cmd.Run(); nil == err {
        return nil
    }

    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("未找到nginx命令，是否安装nginx服务器 [y/N] :")

    input, err := reader.ReadString('\n')
    if nil != err {
        return errors.New(fmt.Sprintf("接收输入失败: %v", err))
    }

    if "y" == strings.Trim(input, "\n") {
        system, err := utils.GetOsType()
        if nil != err {
            return err
        }

        switch system {
        case "debian":
            if err := utils.ExecCommandOutputToLog("apt", []string{"install", "nginx", "-y"}, nil); nil != err {
                return errors.New(fmt.Sprintf("安装nginx失败: %v", err))
            }
        case "redhat":
            if err := utils.ExecCommandOutputToLog("yum", []string{"install", "nginx", "-y"}, nil); nil != err {
                return errors.New(fmt.Sprintf("安装nginx失败: %v", err))
            }
        default:
            return errors.New("不支持的系统类型")
        }

        cmd = exec.Command("nginx", "-v")
        if err := cmd.Run(); nil != err {
            return errors.New(fmt.Sprintf("安装nginx服务器失败: %v", err))
        }

        return nil
    }

    return errors.New("未安装nginx服务器")
}

// generateNginxConfig 生成nginx配置文件
func (c *Cmd) generateNginxConfig(hostname string) error {
    if err := os.MkdirAll("/var/www/html", 0755); nil != err {
        return err
    }

    hostConfig := strings.Replace(template, "{host}", hostname, 1)

    system, err := utils.GetOsType()
    if nil != err {
        return err
    }

    var configPath string

    switch system {
    case "debian":
        if err := os.MkdirAll("/etc/nginx/sites-enabled", 0755); nil != err {
            return err
        }
        configPath = fmt.Sprintf("/etc/nginx/sites-enabled/%v.conf", hostname)
    case "redhat":
        if err := os.MkdirAll("/etc/nginx/conf.d", 0755); nil != err {
            return err
        }
        configPath = fmt.Sprintf("/etc/nginx/conf.d/%v.conf", hostname)
    default:
        return errors.New("不支持的系统类型")
    }

    logger.Infof("生成nginx配置文件: %v", configPath)
    file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if nil != err {
        return err
    }

    if _, err := file.Write([]byte(hostConfig)); nil != err {
        return err
    }

    return nil
}

// download 下载文件
func (c *Cmd) download(from, output string) error {
    res, err := http.Get(from)
    if nil != err {
        return err
    }
    defer res.Body.Close()

    out, err := os.Create(output)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, res.Body)

    return err
}
