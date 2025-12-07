package client

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type TcpConnection struct {
	Conn         net.Conn // 链接信息
	ProtocolName string
	Buffer       [1024]byte // 缓冲区存储读取到的数据 需要动态调整
	NLast        int        // 目前接受了多少字节
	Client       *Client    // 所属的客户端
	Run          bool
}

func MakeClient(c *Client, conn net.Conn) (client *TcpConnection, err error) {
	client = &TcpConnection{
		Conn:         conn,
		ProtocolName: c.Network,
		NLast:        0,
		Run:          true,
		Client:       c,
	}
	return
}

// RemoveClient 移除客户端链接
// 并且调用关闭链接的事件函数
func (tc *TcpConnection) RemoveClient() {
	tc.Client.CallEventFunc("close", tc)
	tc.NLast = 0    // 重置缓冲区位置
	tc.Conn.Close() // 关闭连接 很重要
	tc.Run = false
	// tc.Client.RemoveClient(tc)
}

// HandleMessage 处理消息
func (tc *TcpConnection) HandleMessage() {

	// tc.Run 为 false 的时候就退出连接
	defer wg.Done()

	for tc.Run {
		recvBytes, err := tc.Conn.Read(tc.Buffer[tc.NLast:]) // 累加数据
		if err != nil {
			if err == io.EOF {
				// 对端关闭了连接
				tc.RemoveClient()
				// 不能往下继续执行
				return
			}
		}

		tc.NLast += recvBytes

		switch tc.ProtocolName {
		case "tcp":
			// fmt.Println("receive: ", string(tc.Buffer[0:tc.NLast]))
			tc.Client.CallEventFunc("receive", tc, tc.Buffer[0:tc.NLast])
			tc.NLast = 0 // 重置位置
		case "stream":
		case "http":
		case "websocket":
		default:
		}
	}
}

func (tc *TcpConnection) Send(msg string) {
	fmt.Println("client send msg: ", msg)
	fmt.Println("protocol name: ", tc.ProtocolName)
	switch tc.ProtocolName {
	case "tcp":
		_ = tc.WriteData([]byte(msg))
	case "stream":
	case "http":
	case "websocket":

	}
}

func (tc *TcpConnection) WriteData(data []byte) (err error) {
	length := len(data)
	writeBytes, err := tc.Conn.Write(data)
	// 长度不一致也会有问题
	// 发生数据必须完整
	if err != nil || writeBytes != length {
		if err == nil {
			err = errors.New("数据发生长度不正确")
		}

		// 否则给你关了客户端
		tc.RemoveClient()
	}
	return
}
