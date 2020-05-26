package utils

import (
    "context"
    "golang.org/x/time/rate"
    "io"
    "time"
)

// CopyLimiterWithCtx 基于context的Copy
func CopyLimiterWithCtx(ctx context.Context, reader io.Reader, writer io.Writer, limiter *rate.Limiter) chan error {
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

                nw, ew := writer.Write(buf[0:nr])
                if ew != nil {
                    ch <- ew

                    return
                }

                if nr != nw {
                    ch <- io.ErrShortWrite

                    return
                }

                ch <- nil
            }
        }
    }()

    return ch
}

// CopyWithCtx 基于context的Copy
func CopyWithCtx(ctx context.Context, reader io.Reader, writer io.Writer) chan error {
    ch := make(chan error)

    go func() {
        buf := make([]byte, 1024)

        for {
            select {
            case <-ctx.Done():
                return
            default:
                nr, er := reader.Read(buf)
                if er != nil {
                    ch <- er

                    return
                }

                if nr <= 0 {
                    ch <- io.EOF

                    return
                }

                nw, ew := writer.Write(buf[0:nr])
                if ew != nil {
                    ch <- ew

                    return
                }

                if nr != nw {
                    ch <- io.ErrShortWrite

                    return
                }

                ch <- nil
            }
        }
    }()

    return ch
}
