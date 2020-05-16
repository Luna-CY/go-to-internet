package utils

import (
    "golang.org/x/time/rate"
    "io"
    "time"
)

// Copy 替代io.Copy，支持速率限制
func Copy(writer io.Writer, reader io.Reader, limiter *rate.Limiter) (written int64, err error) {
    if nil == limiter {
        return io.Copy(writer, reader)
    }

    buf := make([]byte, limiter.Burst())
    for {
        if !limiter.AllowN(time.Now(), limiter.Burst()) {
            continue
        }

        nr, er := reader.Read(buf)
        if nr > 0 {
            nw, ew := writer.Write(buf[0:nr])
            if nw > 0 {
                written += int64(nw)
            }
            if ew != nil {
                err = ew
                break
            }
            if nr != nw {
                err = io.ErrShortWrite
                break
            }
        }
        if er != nil {
            if er != io.EOF {
                err = er
            }
            break
        }
    }
    return written, err
}

// Bind 对writer和reader进行绑定
// 返回一个通道，无错误完成时向通道输入0，发生错误时向通道输入1
func Bind(dst io.Writer, src io.Reader, limiter *rate.Limiter) chan int {
    ch := make(chan int)

    go func() {
        code := 0

        if _, err := Copy(dst, src, limiter); nil != err && io.EOF != err {
            code = 1
        }

        ch <- code
    }()

    return ch
}
