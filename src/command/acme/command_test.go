package acme

import "testing"

// 在mac上运行
func TestGetOsType(t *testing.T) {
    system, err := getOsType()
    if nil != err {
        t.Errorf("获取文件系统类型失败: %v\n", err)
    }

    // 测试失败了，居然在mac上找到了apt路径，/Library/Java/JavaVirtualMachines/jdk-12.0.2.jdk/Contents/Home/bin/apt
    if "unknown" != system {
        t.Errorf("测试失败，系统类型为：%v\n", system)
    }
}
