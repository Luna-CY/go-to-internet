package user

import "testing"

func TestConfig_Validate_ListAndAdd(t *testing.T) {
    c := Config{}
    if c.Validate() {
        t.Error("空配置测试失败")
    }

    c = Config{Config: "xxx"}
    if c.Validate() {
        t.Error("空配置测试失败")
    }

    c = Config{Config: "xxx", List: true}
    if !c.Validate() {
        t.Error("列表正确配置测试失败")
    }

    c = Config{Config: "xxx", Add: true}
    if c.Validate() {
        t.Error("添加错误配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "", Password: ""}
    if c.Validate() {
        t.Error("添加错误配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "test", Password: ""}
    if c.Validate() {
        t.Error("添加错误配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "", Password: "password"}
    if c.Validate() {
        t.Error("添加错误配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "test", Password: "password"}
    if !c.Validate() {
        t.Error("添加正确配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "test", Password: "password", MaxRate: -1024}
    if !c.Validate() {
        t.Error("添加正确配置测试失败")
    }

    c = Config{Config: "xxx", Add: true, Username: "test", Password: "password", MaxRate: 1024, MaxConnection: -100}
    if !c.Validate() {
        t.Error("添加正确配置测试失败")
    }
}

func TestConfig_Validate_UpdAndDel(t *testing.T) {
    c := Config{Config: "xxx", Upd: true}
    if c.Validate() {
        t.Error("更新错误配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", MaxRate: -1, MaxConnection: -1}
    if c.Validate() {
        t.Error("更新错误配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", Password: "password"}
    if !c.Validate() {
        t.Error("更新正确配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", MaxRate: -1024, MaxConnection: 100}
    if !c.Validate() {
        t.Error("更新正确配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", Password: "password", MaxRate: 1024, MaxConnection: -100}
    if !c.Validate() {
        t.Error("更新正确配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", Expired: "-"}
    if !c.Validate() {
        t.Error("更新正确配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", Expired: "2020-02-02T23:59:59"}
    if !c.Validate() {
        t.Error("更新正确配置测试失败")
    }

    c = Config{Config: "xxx", Upd: true, Username: "test", Password: "password", Expired: "2020-02-30T23:59:59"}
    if c.Validate() {
        t.Error("更新错误配置测试失败")
    }

    c = Config{Config: "xxx", Del: true}
    if c.Validate() {
        t.Error("删除错误配置测试失败")
    }

    c = Config{Config: "xxx", Del: true, Username: ""}
    if c.Validate() {
        t.Error("删除错误配置测试失败")
    }

    c = Config{Config: "xxx", Del: true, Username: "test"}
    if !c.Validate() {
        t.Error("删除正确配置测试失败")
    }
}
