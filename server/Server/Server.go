package server

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
)

type Server struct {
	Network      string // tcp/udp
	Address      string // ip:port
	Listen       net.Listener
	Clients      map[string]*TcpConnection // 链接的客户端
	MaxClientNum int                       // 最大的客户端连接数的限制
	ClientNum    int                       // 目前的连接数

	OnError   func(err string)
	OnStart   func(srv *Server)
	OnConnect func(srv *Server, client *TcpConnection)
	OnClose   func(srv *Server, client *TcpConnection)
	OnReceive func(srv *Server, client *TcpConnection, data []byte)
}

func (srv *Server) CallEventFunc(eventName string, args ...interface{}) {
	switch eventName {
	case "error":
		srv.OnError(args[0].(string)) // 断言
	case "start":
		srv.OnStart(srv) // 断言
	case "connect":
		srv.OnConnect(srv, args[0].(*TcpConnection))
	case "close":
		srv.OnClose(srv, args[0].(*TcpConnection))
	case "receive":
		srv.OnReceive(srv, args[0].(*TcpConnection), args[1].([]byte))
	default:
		fmt.Println("unknown event name:", eventName)
	}
}

// StartInfo 启动信息
func (srv *Server) StartInfo() {
	fmt.Println("\u001B[33;40mListen success on ", srv.Address, "\u001B[0m")
	fmt.Println("\u001B[33;40mUsing protocol ", srv.Network, "\u001B[0m")

	fmt.Println("\u001B[33;40mplatform  ", runtime.GOOS, "\u001B[0m")
	fmt.Println("\u001B[33;40mcpu Num ", runtime.NumCPU(), "\u001B[0m")
	fmt.Println("\u001B[33;40mversion ", runtime.Version(), "\u001B[0m")
	fmt.Println("\u001B[33;40mPID ", os.Getpid(), "\u001B[0m")
	fmt.Println("------------------------------------------------------------")
}

func (srv *Server) Start() {
	listen, err := net.Listen(srv.Network, srv.Address)
	if err != nil {
		// 出问题了就调用一个错误处理函数
		// srv.OnError(err.Error())
		srv.CallEventFunc("error", err.Error())
		return
	}

	srv.Listen = listen
	defer srv.Listen.Close()

	// 显示一下启动信息
	srv.StartInfo()
	srv.CallEventFunc("start")
	srv.EventLoop()
}

func (srv *Server) AddClient(client *TcpConnection) {
	// ip:port 作为key 确保唯一性
	// fd 会有重复 文件描述符(fd)
	// tcpConn := client.Conn.(*net.TCPConn)
	// f, _ := tcpConn.File()
	// f.Fd() // GOLANG 里面不是唯一,在其他编程语言里是唯一的 底层 fcntl

	// 底层跑的还是多线程
	// 如果出现问题,需要加锁
	srv.Clients[client.Conn.RemoteAddr().String()] = client
	srv.ClientNum++
	// fmt.Println("当前链接数:", len(srv.Clients))
}

func (srv *Server) RemoveClient(client *TcpConnection) {
	ip := client.Conn.RemoteAddr().String()
	_, ok := srv.Clients[ip]
	if ok {
		// 找到就干掉他
		delete(srv.Clients, ip)
		// 客户端数量减少
		srv.ClientNum--
	}
}

func (srv *Server) EventLoop() {
	for {
		conn, err := srv.Listen.Accept()
		if err != nil {
			srv.CallEventFunc("error", err.Error())
			return
		}

		if srv.ClientNum > srv.MaxClientNum {
			srv.CallEventFunc("error", "超过最大连接数限制"+strconv.Itoa(srv.MaxClientNum))
			return
		}

		// 需要将当前服务传递过去
		client, err := MakeClient(srv, conn, srv.Network)
		if err != nil {
			srv.CallEventFunc("error", err.Error())
			return
		}

		// 收到链接之后,将链接存储起来
		srv.AddClient(client)
		srv.CallEventFunc("connect", client)
		go client.HandleMessage()
	}
}
