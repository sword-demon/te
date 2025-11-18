package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	if err != nil {
		fmt.Println("connect err", err)
		return
	}

	// message := "china and american\r\ni like money\r\n"
	// 测试 粘包和少包

	var message []byte = []byte("china and american")
	conn.Write(message[0:5]) // 只发一点点

	time.Sleep(time.Second * 2)

	conn.Write(message[5:])
	conn.Close()
}
