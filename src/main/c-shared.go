package main

// #cgo LDFLAGS: --sysroot=$NDK_ROOT/platforms/android-21/arch-arm
import "C"

//export getHello
func GetHello() *C.char {
    return C.CString("hello")
}

func main() {}
