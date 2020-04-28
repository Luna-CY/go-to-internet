package main

import (
    "gitee.com/Luna-CY/go-to-internet/proxy"
    "net/http"
)

func main() {
    _ = http.ListenAndServe(":4433", &proxy.Proxy{})
}
