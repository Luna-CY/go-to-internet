package main

import "C"

//export hello
func hello() string {
    return "hello from golang"
}

func main() {}
