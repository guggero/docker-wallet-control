package util

import (
    "runtime"
    "log"
)

func LogError(err error) (b bool) {
    if err != nil {
        _, fn, line, _ := runtime.Caller(1)
        log.Printf("[error] %s:%d %v", fn, line, err)
        b = true
    }
    return
}