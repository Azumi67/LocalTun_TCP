package smuxserver

import (
	"github.com/xtaci/smux"
	"github.com/sirupsen/logrus"
	"github.com/Azumi67/LocalTun_TCP/server"
	"net"
)

var log = logrus.New()

func HandleSmux(conn net.Conn, tun *tunserver.Interface, verbose bool) {
	smuxConfig := smux.DefaultConfig()
	session, err := smux.Server(conn, smuxConfig)
	if err != nil {
		log.Warnf("creating Smux session for server failed: %v", err)
		return
	}
	defer session.Close()

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Warnf("Accepting Smux stream failed: %v", err)
			return
		}

		clientToTun := make(chan []byte, 100)
		tunToClient := make(chan []byte, 100)

		go fromClient(stream, tun, clientToTun, verbose)
		go toTun(tun, clientToTun, verbose)
		go fromTun(tun, tunToClient, verbose)
		go toClient(stream, tunToClient, verbose)

		select {}
	}
}

func fromClient(conn net.Conn, tun *tunserver.Interface, clientToTun chan []byte, verbose bool) {
	defer close(clientToTun)

	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			log.Warnf("Couldn't read Smux packet length from client: %v", err)
			return
		}

		length := int(pcktLength[0])<<8 | int(pcktLength[1])
		buff := make([]byte, length)

		_, err := conn.Read(buff)
		if err != nil {
			log.Warnf("Couldn't read Smux data from client: %v", err)
			return
		}

		clientToTun <- buff
	}
}

func toTun(tun *tunserver.Interface, clientToTun chan []byte, verbose bool) {
	for buff := range clientToTun {
		if _, err := tun.Write(buff); err != nil {
			log.Warnf("Couldn't write to TUN device: %v", err)
		}
	}
}

func fromTun(tun *tunserver.Interface, tunToClient chan []byte, verbose bool) {
	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			log.Warnf("Couldn't read from TUN device: %v", err)
			continue
		}

		packet := make([]byte, 2+n)
		packet[0] = byte(n >> 8)
		packet[1] = byte(n)
		copy(packet[2:], buff[:n])

		tunToClient <- packet
	}
}

func toClient(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			log.Warnf("Couldn't write to client: %v", err)
		}
	}
}
