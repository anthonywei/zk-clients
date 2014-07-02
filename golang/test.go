package main

import (
    //"./zk"
    "./mtconfig"
    "fmt"
    "time"
)

func main(){
    mc, err := mtconfig.NewMtConfigServer("")

    if(err != nil) {
        panic(err)
    }

    value, _ := mc.Get("inftest", "common", "test")

    fmt.Printf("%s\n", value)

    for {
        value, _ = mc.Get("inftest", "common", "test")
        fmt.Printf("%s\n", value)
        time.Sleep(time.Second * 1)
        //TODO your own businsess
    }
}
