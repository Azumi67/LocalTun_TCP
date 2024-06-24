package tcpnodelayclient

import (
	"net"
	"syscall"
	"time"
)

func noDelay(conn *net.TCPConn) error {
	err := conn.noDelay(true)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return err
	}

	if err := syscall.SetsockoptInt(conn.SocketConn().Sysfd, syscall.IPPROTO_TCP, syscall.TCP_FASTOPEN, 1); err != nil {
		return err
	}

	return nil
}
