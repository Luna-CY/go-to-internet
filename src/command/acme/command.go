package acme

import (
    "bufio"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
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

        if err := execCommand(command, []string{"--issue", "-d", config.Hostname, "--nginx"}, nil, true); nil != err {
            return err
        }

        logger.Info("申请证书完成")
    }

    return nil
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
