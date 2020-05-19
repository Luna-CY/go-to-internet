package acme

import "testing"

func TestConfig_Validate(t *testing.T) {
    c := Config{}
    if c.Validate() {
        t.Error("空配置测试失败")
    }

    c = Config{Install: true}
    if !c.Validate() {
        t.Error("安装配置测试失败")
    }

    c = Config{Issue: true}
    if c.Validate() {
        t.Error("申请证书错误配置测试失败")
    }

    c = Config{Issue: true, Hostname: ""}
    if c.Validate() {
        t.Error("申请证书空域名配置测试失败")
    }

    c = Config{Issue: true, Hostname: "www.example.com"}
    if !c.Validate() {
        t.Error("申请证书正确配置测试失败")
    }

    c = Config{Issue: true, Hostname: "www.example.com", Nginx: true}
    if !c.Validate() {
        t.Error("申请证书nginx模式配置测试失败")
    }

    c = Config{Issue: true, Hostname: "www.example.com", Standalone: true}
    if !c.Validate() {
        t.Error("申请证书standalone模式配置测试失败")
    }
}
