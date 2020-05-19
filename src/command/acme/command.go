package acme

import (
    "bufio"
    "errors"
    "fmt"
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

const AcmePath = "https://raw.githubusercontent.com/acmesh-official/acme.sh/master/acme.sh"

// Exec 执行acme子命令
func Exec(config *Config) error {
    switch {
    case config.Install:
        home, err := os.UserHomeDir()
        if nil != err {
            return err
        }
        acmePath := path.Join(home, ".acme.sh")

        info, err := os.Stat(acmePath)
        if nil == err {
            if !info.IsDir() {
                return errors.New(fmt.Sprintf("[%v]路径已存在并且不是一个目录", acmePath))
            }

            logger.Infof("已安装acme.sh工具，重新安装请删除[%v]目录后重新执行", acmePath)

            return nil
        }

        if !os.IsNotExist(err) {
            return errors.New(fmt.Sprintf("获取路径信息失败: %v", err))
        }

        output := path.Join(os.TempDir(), "install-acme.sh")

        if err := download(AcmePath, output); nil != err {
            return err
        }

        if err := os.Chmod(output, os.FileMode(0755)); nil != err {
            return err
        }

        if err := execCommand("sh", []string{"-c", output}, &[]string{"INSTALLONLINE=1"}, true); nil != err {
            return err
        }

        if err := unix.Unlink(output); nil != err {
            return errors.New(fmt.Sprintf("删除安装脚本失败: %v", err))
        }

        logger.Info("安装完成")
    case config.Issue:
        home, err := os.UserHomeDir()
        if nil != err {
            return err
        }
        command := path.Join(home, ".acme.sh", "acme.sh")
        info, err := os.Stat(command)
        if nil != err {
            return err
        }

        if info.IsDir() {
            return errors.New("无法找到acme.sh工具")
        }

        switch {
        case config.Nginx:
            if err := generateNginxConfig(config.Hostname); nil != err {
                return errors.New(fmt.Sprintf("创建nginx配置失败: %v", err))
            }

            if err := execCommand(command, []string{"--issue", "-d", config.Hostname, "--nginx"}, nil, true); nil != err {
                return err
            }
        default:
            if err := execCommand(command, []string{"--issue", "-d", config.Hostname, "--standalone"}, nil, true); nil != err {
                return err
            }
        }

        logger.Info("申请证书完成")
    }

    return nil
}

// generateNginxConfig 生成nginx配置文件
func generateNginxConfig(hostname string) error {
    hostConfig := strings.Replace(template, "{host}", hostname, 1)

    system, err := getOsType()
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
    case "centos":
        if err := os.MkdirAll("/etc/nginx/conf.d", 0755); nil != err {
            return err
        }
        configPath = fmt.Sprintf("/etc/nginx/conf.d/%v.conf", hostname)
    default:
        return errors.New("不支持的系统类型")
    }

    file, err := os.OpenFile(configPath, os.O_WRONLY, 0644)
    if nil != err {
        return err
    }

    if _, err := file.Write([]byte(hostConfig)); nil != err {
        return err
    }

    return nil
}

// getOsType 获取文件系统类型
func getOsType() (string, error) {

    if debian, err := utils.FileExists("/usr/bin/apt"); nil != err || debian {
        if debian {
            return "debian", nil
        }

        return "", err
    }

    if centos, err := utils.FileExists("/usr/bin/yum"); nil != err || centos {
        if centos {
            return "centos", nil
        }

        return "", err
    }

    return "unknown", nil
}

// execCommand 执行命令
func execCommand(name string, args []string, env *[]string, output bool) error {
    cmd := exec.Command(name, args...)
    cmd.Env = os.Environ()

    if nil != env {
        for _, value := range *env {
            cmd.Env = append(cmd.Env, value)
        }
    }

    if output {
        stdout, err := cmd.StdoutPipe()
        if nil != err {
            return err
        }
        defer stdout.Close()

        if err := cmd.Start(); nil != err {
            return err
        }

        reader := bufio.NewReader(stdout)
        for {
            line, err := reader.ReadString('\n')
            if err != nil || io.EOF == err {
                break
            }
            logger.Info(strings.Trim(line, "\n"))
        }

        if err := cmd.Wait(); nil != err {
            return err
        }
    } else {
        if err := cmd.Run(); nil != err {
            return err
        }
    }

    return nil
}

// download 下载文件
func download(from, output string) error {
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
