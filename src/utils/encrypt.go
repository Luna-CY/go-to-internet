package utils

import (
    "crypto/md5"
    "fmt"
)

// EncryptPassword 计算加密后的密码
func EncryptPassword(password string) string {
    passwordBytes := md5.Sum([]byte(password))

    newPasswordBytes := make([]byte, 8)
    for i, b := range passwordBytes {
        if 0 == i%2 {
            newPasswordBytes[i/2] = b
        }
    }

    return fmt.Sprintf("%x", newPasswordBytes)
}
