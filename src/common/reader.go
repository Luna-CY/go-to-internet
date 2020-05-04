package common

import (
    "fmt"
    "io"
)

// ReadAll 替代ioutil.ReadAll方法
func ReadAll(reader io.Reader) ([]byte, error) {
    var data []byte
    buffer := make([]byte, 256)

    for {
        n, err := reader.Read(buffer)
        data = append(data, buffer[:n]...)

        if nil != err && io.EOF != err {
            fmt.Println(err)

            return nil, err
        }

        if io.EOF == err || len(buffer) > n {
            return data, nil
        }
    }
}

// Copy 替代io.Copy方法
func Copy(writer io.Writer, reader io.Reader) (int, error) {
    buffer := make([]byte, 256)
    counter := 0

    for {
        n, err := reader.Read(buffer)
        counter += n
        _, _ = writer.Write(buffer[:n])

        if nil != err || len(buffer) > n {
            return counter, err
        }
    }
}
