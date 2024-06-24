package tcpnodelayclient

import (
	"net"
)

func NoDelay(conn net.Conn, tcpNoDelay bool) error {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		return tcpConn.SetNoDelay(tcpNoDelay)
	}
	return nil
}
