package main

import (
	"fmt"

	server "github.com/sword-demon/te/server/Server"
)

func main() {
	(&server.Server{
		Network:      "tcp", // stream, http, websocket
		Address:      "0.0.0.0:9501",
		MaxClientNum: 1024 * 100,
		ClientNum:    0,
		Clients:      make(map[string]*server.TcpConnection), // 分配好空间
		OnConnect: func(srv *server.Server, client *server.TcpConnection) {
			fmt.Println("有新的链接进来了:", client.Conn.RemoteAddr().String())
		},
		OnError: func(err string) {
			fmt.Println("[OnError] = ", err)
		},
		OnStart: func(srv *server.Server) {
			fmt.Println("成功启动服务")
		},
	}).Start()
}
