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
    "runtime"
    "strings"
)

// Exec 执行acme子命令
func Exec(config *Config) error {
    switch {
    case config.Install:
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

        if err := download(common.AcmePath, output); nil != err {
            logger.Errorf("下载安装脚本失败: %v", err)

            break
        }

        if err := os.Chmod(output, os.FileMode(0755)); nil != err {
            logger.Errorf("修改文件权限失败: %v", err)

            break
        }

        if err := execCommand("sh", []string{"-c", output}, &[]string{"INSTALLONLINE=1"}, true); nil != err {
            logger.Errorf("安装acme.sh失败: %v", err)

            break
        }

        if err := unix.Unlink(output); nil != err {
            logger.Errorf("删除安装脚本失败: %v", err)

            break
        }

        logger.Info("安装完成")
    case config.Issue:
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
        case config.Nginx:
            if err := checkAndInstallNginx(); nil != err {
                logger.Error(err)

                break
            }

            if err := generateNginxConfig(config.Hostname); nil != err {
                logger.Errorf("创建nginx配置失败: %v", err)

                break
            }

            isExist, err := utils.FileExists(fmt.Sprintf("/root/.acme.sh/%v", config.Hostname))
            if nil != err {
                logger.Errorf("检查证书路径失败: %v", err)

                break
            }

            if isExist {
                logger.Info("该域名已存在证书")

                break
            }

            if err := execCommand(command, []string{"--issue", "-d", config.Hostname, "--nginx"}, nil, true); nil != err {
                logger.Errorf("申请证书失败: %v", err)

                break
            }

            logger.Info("申请证书完成")
        default:
            if err := execCommand(command, []string{"--issue", "-d", config.Hostname, "--standalone"}, nil, true); nil != err {
                logger.Errorf("申请证书失败: %v", err)

                break
            }

            logger.Info("申请证书完成")
        }
    }

    return nil
}

// checkAndInstallNginx 检查nginx是否存在，不存在时安装nginx
func checkAndInstallNginx() error {
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
        system, err := getOsType()
        if nil != err {
            return err
        }

        switch system {
        case "debian":
            if err := execCommand("apt", []string{"install", "nginx", "-y"}, nil, true); nil != err {
                return errors.New(fmt.Sprintf("安装nginx失败: %v", err))
            }
        case "redhat":
            if err := execCommand("yum", []string{"install", "nginx", "-y"}, nil, true); nil != err {
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
func generateNginxConfig(hostname string) error {
    if err := os.MkdirAll("/var/www/html", 0755); nil != err {
        return err
    }

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
    case "redhat":
        if err := os.MkdirAll("/etc/nginx/conf.d", 0755); nil != err {
            return err
        }
        configPath = fmt.Sprintf("/etc/nginx/conf.d/%v.conf", hostname)
    default:
        return errors.New("不支持的系统类型")
    }

    file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
    switch runtime.GOOS {
    case "darwin":
        return runtime.GOOS, nil
    case "linux":
        // ubuntu视为debian
        if debian, err := utils.FileExists("/etc/debian_version"); nil != err || debian {
            if debian {
                return "debian", nil
            }

            return "", err
        }

        // centos/fedora视为redhat
        if redhat, err := utils.FileExists("/etc/redhat-version"); nil != err || redhat {
            if redhat {
                return "redhat", nil
            }

            return "", err
        }

        return "unknown", nil
    case "windows":
        return runtime.GOOS, nil
    default:
        return "unknown", nil
    }
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

        stderr, err := cmd.StderrPipe()
        if nil != err {
            return err
        }
        defer stderr.Close()

        if err := cmd.Start(); nil != err {
            return err
        }

        outReader := bufio.NewReader(stdout)
        for {
            line, err := outReader.ReadString('\n')
            if err != nil || io.EOF == err {
                break
            }
            logger.Info(strings.Trim(line, "\n"))
        }

        errReader := bufio.NewReader(stderr)
        for {
            line, err := errReader.ReadString('\n')
            if err != nil || io.EOF == err {
                break
            }
            logger.Error(strings.Trim(line, "\n"))
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
