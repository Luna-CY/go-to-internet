package acme

import "testing"

// 在mac上运行
func TestGetOsType(t *testing.T) {
    system, err := getOsType()
    if nil != err {
        t.Errorf("获取文件系统类型失败: %v\n", err)
    }

    if "darwin" != system {
        t.Errorf("测试失败，系统类型为：%v\n", system)
    }
}
