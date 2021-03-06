package utils

import (
    "bufio"
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "io"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

// FileExists 检查文件路径是否存在
func FileExists(filepath string) (bool, error) {
    if _, err := os.Stat(filepath); nil != err {
        if os.IsNotExist(err) {
            return false, nil
        }

        return false, err
    }

    return true, nil
}

// ExecCommandGetStdout 执行命令并获取stdout/stderr
func ExecCommandGetStdout(name string, args []string, env *[]string) ([]string, error) {
    cmd := exec.Command(name, args...)
    cmd.Env = os.Environ()

    if nil != env {
        for _, value := range *env {
            cmd.Env = append(cmd.Env, value)
        }
    }

    stdout, err := cmd.StdoutPipe()
    if nil != err {
        return nil, err
    }
    defer stdout.Close()

    stderr, err := cmd.StderrPipe()
    if nil != err {
        return nil, err
    }
    defer stderr.Close()

    if err := cmd.Start(); nil != err {
        return nil, err
    }

    result := make([]string, 1)
    outReader := bufio.NewReader(stdout)
    for {
        line, err := outReader.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        result = append(result, strings.Trim(line, "\n"))
    }

    errReader := bufio.NewReader(stderr)
    for {
        line, err := errReader.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        result = append(result, strings.Trim(line, "\n"))
    }

    if err := cmd.Wait(); nil != err {
        return nil, err
    }

    return result, nil
}

// ExecCommandOutputToLog 执行命令并将输入打印到日志
func ExecCommandOutputToLog(name string, args []string, env *[]string) error {
    cmd := exec.Command(name, args...)
    cmd.Env = os.Environ()

    if nil != env {
        for _, value := range *env {
            cmd.Env = append(cmd.Env, value)
        }
    }

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
        logger.Info(strings.Trim(line, "\n"))
    }

    if err := cmd.Wait(); nil != err {
        return err
    }

    return nil
}

// GetOsType 获取文件系统类型：debian/redhat/windows/darwin/unknown
func GetOsType() (string, error) {
    switch runtime.GOOS {
    case "darwin":
        return runtime.GOOS, nil
    case "linux":
        // ubuntu视为debian
        if debian, err := FileExists("/etc/debian_version"); nil != err || debian {
            if debian {
                return "debian", nil
            }

            return "", err
        }

        // centos/fedora视为redhat
        if redhat, err := FileExists("/etc/redhat-version"); nil != err || redhat {
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
