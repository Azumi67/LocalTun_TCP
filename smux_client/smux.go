package smuxclient

import (
	"net"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
)

var log = logrus.New()

func HandleSmux(conn net.Conn, tun *water.Interface, verbose bool) {
	smuxConfig := smux.DefaultConfig()
	session, err := smux.Client(conn, smuxConfig)
	if err != nil {
		log.Warnf("Creating smux client session failed: %v", err)
		return
	}
	defer session.Close()

	stream, err := session.OpenStream()
	if err != nil {
		log.Warnf("Opening smux stream failed: %v", err)
		return
	}
	defer stream.Close()

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromServer(stream, tun, clientToTun, verbose)
	go tunclient.ToTun(tun, clientToTun, verbose)
	go tunclient.FromTun(tun, tunToClient, verbose)
	go toServer(stream, tunToClient, verbose)

	select {}
}

func fromServer(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			log.Warnf("Couldn't read packet length from server: %v", err)
			close(clientToTun)
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := tunclient.Data(conn, buff)
		if err != nil {
			log.Warnf("Couldn't read data from server: %v", err)
			close(clientToTun)
			return
		}

		clientToTun <- buff
	}
}

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		_, err := conn.Write(packet)
		if err != nil {
			log.Warnf("Couldn't write to server: %v", err)
		}
	}
}
