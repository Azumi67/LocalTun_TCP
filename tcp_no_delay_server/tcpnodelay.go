package tcpnodelayserver

import (
	"net"
	"github.com/songgao/water"
	"sync"
	"github.com/sirupsen/logrus"
	"github.com/Azumi67/LocalTun_TCP/server"
)

var log = logrus.New()

func HandleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, tcpNoDelay bool) {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(tcpNoDelay)
	}

	buff := make([]byte, len(secretKey))
	if _, err := conn.Read(buff); err != nil {
		log.Warnf("Couldn't read authentication key: %v", err)
		conn.Close()
		return
	}

	if string(buff) != secretKey {
		log.Warn("Wrong authentication key")
		conn.Close()
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tunserver.FromClient(conn, tun, clientToTun, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tunserver.ToTun(tun, clientToTun, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tunserver.FromTun(tun, tunToClient, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tunserver.ToClient(conn, tunToClient, verbose)
	}()

	wg.Wait()
}
