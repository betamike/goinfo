package main

import (
    "github.com/betamike/goinfo"
)

func main() {
    err := goinfo.Start("/opt/go/example")
    if err != nil {
        panic("Oh no! could not mount our monitoring file system: " + err.Error())
    }

    //Start app business logic
    select{}
}
