package main

import (
    "os"
    "fmt"
    "net"
    "time"
//  "io"
)

func main() {
    conn, err := net.Dial("udp", os.Args[1])
    defer conn.Close()
    if err != nil {
        os.Exit(1)
    }
    for {
      conn.Write([]byte("Hello world!"))
      fmt.Println("send msg")
      time.Sleep(time.Second)
    }
    //conn.Write([]byte("Hello world!"))

    //fmt.Println("send msg")

    var msg [20]byte
    conn.Read(msg[0:])

    fmt.Println("msg is", string(msg[0:10]))
}
