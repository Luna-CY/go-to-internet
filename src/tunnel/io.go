package tunnel

import (
    "context"
    "golang.org/x/time/rate"
    "io"
    "net"
    "time"
)

// CopyLimiterWithCtxToMessageProtocol 基于context的Copy
func CopyLimiterWithCtxToMessageProtocol(ctx context.Context, reader net.Conn, writer net.Conn, limiter *rate.Limiter) chan error {
    ch := make(chan error)

    go func() {
        var length int

        if nil != limiter {
            length = limiter.Burst()
        } else {
            length = 1024
        }

        buf := make([]byte, length)

        for {
            select {
            case <-ctx.Done():
                return
            default:
                if nil != limiter {
                    if !limiter.AllowN(time.Now(), limiter.Burst()) {
                        ch <- nil

                        continue
                    }
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

                message := NewDataMessage(writer, buf[0:nr])
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
func CopyWithCtxFromMessageProtocol(ctx context.Context, reader net.Conn, writer net.Conn) (chan error, chan *MessageProtocol) {
    ch := make(chan error)
    chm := make(chan *MessageProtocol)

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                message := NewEmptyMessage(reader)
                if er := message.Receive(); nil != er {
                    ch <- er

                    return
                }

                // 如果不是数据指令，那么返回这个消息
                if CmdData != message.Cmd {
                    chm <- message

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

    return ch, chm
}
