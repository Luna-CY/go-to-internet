package utils

import "os"

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
