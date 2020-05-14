package acme

import (
    "bufio"
    "errors"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "io"
    "os"
    "os/exec"
    "path"
    "strings"
)

// Exec 执行acme子命令
func Exec(config *Config) error {
    switch {
    case config.Install:
        output := path.Join(os.TempDir(), "install-acme.sh")

        cmd := exec.Command("curl", "https://get.acme.sh", "-o", output)
        if err := cmd.Run(); nil != err {
            return err
        }

        if err := os.Chmod(output, os.FileMode(0755)); nil != err {
            return err
        }

        if err := execCmd(output); nil != err {
            return err
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

        if err := execCmd(command, "--issue", "-d", config.Hostname, "--nginx"); nil != err {
            return err
        }

        logger.Info("申请证书完成")
    }

    return nil
}

// execCmd 执行外部命令并输出结果到logger
func execCmd(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    logger.Infof("执行命令: %v", strings.Join(cmd.Args, " "))

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
        logger.Info(line)
    }

    if err := cmd.Wait(); nil != err {
        return err
    }
    return nil
}
