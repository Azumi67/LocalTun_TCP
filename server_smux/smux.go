package serversmux

import (
	"net"
	"sync"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
	"github.com/sirupsen/logrus"
	"github.com/Azumi67/LocalTun_TCP/server"
)

var log = logrus.New()

func HandleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool) {
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

	handleSmux(conn, tun, verbose)
}

func handleSmux(conn net.Conn, tun *water.Interface, verbose bool) {
	smuxConfig := smux.DefaultConfig()
	session, err := smux.Server(conn, smuxConfig)
	if err != nil {
		log.Warnf("Creating smux for server failed: %v", err)
		conn.Close()
		return
	}
	defer session.Close()

	stream, err := session.AcceptStream()
	if err != nil {
		log.Warnf("Accepting smux stream failed: %v", err)
		conn.Close()
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tunserver.FromClient(stream, tun, clientToTun, verbose)
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
		tunserver.ToClient(stream, tunToClient, verbose)
	}()

	wg.Wait()
}
