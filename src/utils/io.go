package utils

import (
    "context"
    "errors"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/src/tunnel"
    "golang.org/x/time/rate"
    "io"
    "net"
    "time"
)

// CopyLimiterWithCtxToMessageProtocol 基于context的Copy
func CopyLimiterWithCtxToMessageProtocol(ctx context.Context, reader net.Conn, writer net.Conn, limiter *rate.Limiter) chan error {
    ch := make(chan error)

    go func() {
        buf := make([]byte, limiter.Burst())

        for {
            select {
            case <-ctx.Done():
                return
            default:
                if !limiter.AllowN(time.Now(), limiter.Burst()) {
                    ch <- nil

                    continue
                }

                nr, er := reader.Read(buf)
                if er != nil {
                    ch <- er

                    return
                }

                if nr <= 0 {
                    ch <- io.EOF

                    return
                }

                message := tunnel.NewDataMessage(writer, buf[0:nr])
                if ew := message.Send(); nil != ew {
                    ch <- ew

                    return
                }

                ch <- nil
            }
        }
    }()

    return ch
}

// CopyWithCtxFromMessageProtocol 基于context的Copy
func CopyWithCtxFromMessageProtocol(ctx context.Context, reader net.Conn, writer net.Conn) chan error {
    ch := make(chan error)

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                message := tunnel.NewEmptyMessage(reader)
                if er := message.Receive(); nil != er {
                    ch <- er

                    return
                }

                if tunnel.CmdData != message.Cmd || tunnel.MessageCodeSuccess != message.Code {
                    ch <- errors.New(fmt.Sprintf("接收连接消息失败. 响应指令: %v 响应码: %v", message.Cmd, message.Code))

                    return
                }

                nw, ew := writer.Write(message.Data)
                if ew != nil {
                    ch <- ew

                    return
                }

                if len(message.Data) != nw {
                    ch <- io.ErrShortWrite

                    return
                }

                ch <- nil
            }
        }
    }()

    return ch
}
