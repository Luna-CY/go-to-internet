package utils

import (
    "bufio"
    "io"
    "os"
    "os/exec"
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
