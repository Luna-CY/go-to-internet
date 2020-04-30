package common

import "io"

// ReadAll 替代ioutil.ReadAll方法
func ReadAll(reader io.Reader) ([]byte, error) {
    var data []byte
    buffer := make([]byte, 256)

    for {
        n, err := reader.Read(buffer)
        data = append(data, buffer[:n]...)

        if nil != err && io.EOF != err {
            return nil, err
        }

        if io.EOF == err || len(buffer) > n {
            return data, nil
        }
    }
}
