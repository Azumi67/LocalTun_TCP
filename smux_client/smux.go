package smuxclient

import (
	"encoding/binary"
	"net"
	"github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
	"github.com/songgao/water"
)

var log = logrus.New()

func handleSmux(conn net.Conn, tun *water.Interface, verbose bool) {
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
	go toTun(tun, clientToTun, verbose)
	go fromTun(tun, tunToClient, verbose)
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

		_, err := data(conn, buff)
		if err != nil {
			log.Warnf("Couldn't read data from server: %v", err)
			close(clientToTun)
			return
		}

		clientToTun <- buff
	}
}

func toTun(tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for buff := range clientToTun {
		if _, err := tun.Write(buff); err != nil {
			log.Warnf("Couldn't write to TUN device: %v", err)
		}
	}
}

func fromTun(tun *water.Interface, tunToClient chan []byte, verbose bool) {
	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			log.Warnf("Couldn't read from TUN device: %v", err)
			continue
		}

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		tunToClient <- packet
	}
}

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			log.Warnf("Couldn't write to server: %v", err)
		}
	}
}

func data(conn net.Conn, buff []byte) (int, error) {
	totalRead := 0
	for totalRead < len(buff) {
		n, err := conn.Read(buff[totalRead:])
		if err != nil {
			return totalRead, err
		}
		totalRead += n
	}
	return totalRead, nil
}
