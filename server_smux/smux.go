package serversmux

import (
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
	"net"
	"github.com/Azumi67/LocalTun_TCP/server"
)

var log = logrus.New()

func HandleSmux(conn net.Conn, tun *water.Interface, verbose bool) {
	smuxConfig := smux.DefaultConfig()
	session, err := smux.Server(conn, smuxConfig)
	if err != nil {
		log.Warnf("Creating smux server session failed: %v", err)
		return
	}
	defer session.Close()

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Warnf("Accepting smux stream failed: %v", err)
			return
		}

		clientToTun := make(chan []byte, 100)
		tunToClient := make(chan []byte, 100)

		go tunserver.FromClient(stream, tun, clientToTun, verbose)
		go tunserver.ToTun(tun, clientToTun, verbose)
		go tunserver.FromTun(tun, tunToClient, verbose)
		go tunserver.ToClient(stream, tunToClient, verbose)

		select {}
	}
}
