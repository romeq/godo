package utils

import "log"

func Check(err error) {
    if err != nil {
        log.Panicf("error: %v", err)
    }
}

func Ensure(some any) {
    if some == nil {
        log.Panicf("error: %s", some)
    }
}


