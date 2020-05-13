package acme

import "os/exec"

func Exec(config *Config) error {
    if config.Install {
        cmd := exec.Command("curl", "https://get.acme.sh", ">", "/tmp/install-acme.sh")

        if err := cmd.Run(); nil != err {
            return err
        }
    }
    return nil
}
