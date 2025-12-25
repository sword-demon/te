// / 封装一个链接
package server

import (
	"errors"
	"io"
	"net"
)

type TcpConnection struct {
	Conn         net.Conn // 链接信息
	ProtocolName string
	Buffer       [1024]byte // 缓冲区存储读取到的数据 需要动态调整
	NLast        int        // 目前接受了多少字节
	Server       *Server    // 所属的服务
	Run          bool
}

func MakeClient(server *Server, conn net.Conn, protocolName string) (client *TcpConnection, err error) {
	client = &TcpConnection{
		Conn:         conn,
		ProtocolName: protocolName,
		NLast:        0,
		Run:          true,
		Server:       server, // 关联上服务
	}
	return
}

// RemoveClient 移除客户端链接
// 并且调用关闭链接的事件函数
func (tc *TcpConnection) RemoveClient() {
	tc.Server.CallEventFunc("close", tc)
	tc.NLast = 0 // 重置缓冲区位置
	err := tc.Conn.Close()
	if err != nil {
		return
	} // 关闭连接 很重要
	tc.Run = false
	tc.Server.RemoveClient(tc)
}

// HandleMessage 处理消息
func (tc *TcpConnection) HandleMessage() {

	for tc.Run {
		// read recv recvfrom recvmsg 底层函数
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
			tc.Server.CallEventFunc("receive", tc, tc.Buffer[0:tc.NLast])
			// tc.Buffer 这里是覆盖接收
			tc.NLast = 0 // 重置位置
		case "stream":
		case "http":
		case "websocket":
		default:
		}
	}
}

func (tc *TcpConnection) Send(msg string) {
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
