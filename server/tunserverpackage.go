package tunserver

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"net"
	"sync"
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

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fromClient(conn, tun, clientToTun, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		toTun(tun, clientToTun, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fromTun(tun, tunToClient, verbose)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		toClient(conn, tunToClient, verbose)
	}()

	wg.Wait()
}

func fromClient(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			log.Warnf("Couldn't read packet length from client: %v", err)
			close(clientToTun) 
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := data(conn, buff)
		if err != nil {
			log.Warnf("Couldn't read data from client: %v", err)
			close(clientToTun) 
			return
		}

		clientToTun <- buff
	}
}

func toTun(tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		buff, ok := <-clientToTun
		if !ok {
			return 
		}
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

func toClient(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for {
		packet, ok := <-tunToClient
		if !ok {
			return 
		}
		if _, err := conn.Write(packet); err != nil {
			log.Warnf("Couldn't write to client: %v", err)
			return
		}
	}
}

func data(conn net.Conn, buff []byte) (int, error) {
	n, err := conn.Read(buff)
	if err != nil {
		return 0, err
	}
	return n, nil
}
