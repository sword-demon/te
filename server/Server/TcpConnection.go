// / 封装一个链接
package server

import "net"

type TcpConnection struct {
	Conn         net.Conn // 链接信息
	ProtocolName string
}

func MakeClient(conn net.Conn, protocolName string) (client *TcpConnection, err error) {
	client = &TcpConnection{
		Conn:         conn,
		ProtocolName: protocolName,
	}
	return
}

func (tc *TcpConnection) HandleMessage() {
	// for {
	// 	tc.Conn.Read()
	// }
}
