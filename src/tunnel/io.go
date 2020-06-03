package tunnel

import (
    "context"
    "golang.org/x/time/rate"
    "io"
    "net"
)

// CopyLimiterWithCtxToMessageProtocol 基于context的Copy
func CopyLimiterWithCtxToMessageProtocol(ctx context.Context, reader net.Conn, writer net.Conn, limiter *rate.Limiter) chan error {
    ch := make(chan error)

    go func() {
        var length int

        if nil != limiter {
            length = 1024
        } else {
            length = 1024
        }

        buf := make([]byte, length)

        for {
            if nil != limiter {
                if err := limiter.Wait(ctx); nil != err {
                    ch <- err

                    return
                }
            }

            nr, er := reader.Read(buf)
            if nil != ctx.Err() {
                return
            }

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
    }()

    return ch
}

// CopyFromMessageProtocol 基于context的Copy
func CopyFromMessageProtocol(reader net.Conn, writer net.Conn) (chan error, chan *MessageProtocol) {
    ch := make(chan error)
    chm := make(chan *MessageProtocol)

    go func() {
        for {
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
    }()

    return ch, chm
}
