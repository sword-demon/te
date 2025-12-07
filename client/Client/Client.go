package client

import (
	"net"
	"sync"
)

var wg sync.WaitGroup

type Client struct {
	Network string
	Address string
	Connect *TcpConnection

	OnError   func(err string)
	OnStart   func(c *Client)
	OnConnect func(c *Client, conn *TcpConnection)
	OnClose   func(c *Client, conn *TcpConnection)
	OnReceive func(c *Client, conn *TcpConnection, data []byte)
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp4", c.Address)
	if err != nil {
		c.CallEventFunc("error", err.Error())
		return
	}

	c.CallEventFunc("start")
	client, _ := MakeClient(c, conn)
	c.Connect = client
	// 创建客户端连接
	c.CallEventFunc("connect", c.Connect)

	// 多协程使用,当主协程结束时,整个进程就结束了,其他没有结束的子协程,会被强制结束
	wg.Add(1)
	go c.Connect.HandleMessage()
	wg.Wait() // 阻塞当前主线程,主协程
}

func (c *Client) CallEventFunc(eventName string, args ...interface{}) {
	switch eventName {
	case "error":
		c.OnError(args[0].(string)) // 断言
	case "start":
		c.OnStart(c)
	case "connect":
		c.OnConnect(c, args[0].(*TcpConnection)) // 断言
	case "close":
		c.OnClose(c, args[0].(*TcpConnection)) // 断言
	case "receive":
		c.OnReceive(c, args[0].(*TcpConnection), args[1].([]byte)) // 断言
	}
}
