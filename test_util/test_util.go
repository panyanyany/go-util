package test_util

import (
    "os"
    "strings"
)

func IsTesting() bool {
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "-test.v=") {
            return true
        }
    }
    return false
}

func IsTesting2() bool {
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "-test.") {
            return true
        }
    }
    return false
}
