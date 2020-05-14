package acme

import (
    "gitee.com/Luna-CY/go-to-internet/src/logger"
    "golang.org/x/sys/unix"
    "os"
    "os/exec"
    "path"
)

// Exec 执行acme子命令
func Exec(config *Config) error {
    if config.Install {
        output := path.Join(os.TempDir(), "install-acme.sh")

        cmd := exec.Command("curl", "https://get.acme.sh", "-o", output)
        if err := cmd.Run(); nil != err {
            return err
        }

        if err := os.Chmod(output, os.FileMode(0755)); nil != err {
            return err
        }

        cmd = exec.Command(output)
        if err := cmd.Run(); nil != err {
            return err
        }

        if err := unix.Unlink(output); nil != err {
            return err
        }

        logger.Info("安装完成")
    }

    return nil
}
