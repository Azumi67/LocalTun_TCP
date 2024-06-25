package tcpnodelayserver

import (
	"net"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func noDelay(conn net.Conn) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		if err := tcpConn.SetNoDelay(true); err != nil {
			log.Warnf("Setting TCP_NODELAY failed: %v", err)
		}
	}
}
