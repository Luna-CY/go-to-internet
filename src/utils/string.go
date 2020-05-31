package utils

import (
    "math/rand"
)

const charters = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

// RandomString 随机字符串
func RandomString(n int) string {
    b := make([]byte, n)

    for i := range b {
        b[i] = charters[rand.Intn(len(charters))]
    }

    return string(b)
}
