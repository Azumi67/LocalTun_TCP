package tcpnodelay

import (
	"net"
)

func noDelay(conn net.Conn) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Warnf("Setting up TCP nodelay failed")
		return
	}

	err := tcpConn.noDelay(true)
	if err != nil {
		log.Warnf("failed to set up TCP nodelay: %v", err)
	}
}
