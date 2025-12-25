package main

import (
	"fmt"

	client "github.com/sword-demon/te/client/Client"
)

func main() {
	// conn, err := net.Dial("tcp4", "127.0.0.1:9501")
	// if err != nil {
	// 	fmt.Println("connect err", err)
	// 	return
	// }

	// // message := "china and American\r\ni like money\r\n"
	// // 测试 粘包和少包

	// var message []byte = []byte("China and American")
	// conn.Write(message[0:5]) // 只发一点点

	// time.Sleep(time.Second * 2)

	// conn.Write(message[5:])
	// conn.Close()

	// go build -o te_client client/main/main.go
	(&client.Client{
		Network: "tcp",
		Address: "127.0.0.1:9501",
		OnError: func(err string) {
			fmt.Println("客户端服务器启动失败: err", err)
		},
		OnStart: func(c *client.Client) {
			fmt.Println("客户端服务启动成功")
		},
		OnConnect: func(c *client.Client, conn *client.TcpConnection) {
			fmt.Println("成功连接到服务器")
			conn.Send("i am client")
		},
		OnClose: func(c *client.Client, conn *client.TcpConnection) {
			fmt.Println("服务器关闭了连接")
		},
		OnReceive: func(c *client.Client, conn *client.TcpConnection, data []byte) {
			fmt.Println("收到服务器发来的数据: ", string(data))
		},
	}).Start()
}
