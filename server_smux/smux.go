package smuxserver

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/xtaci/smux"
	"github.com/songgao/water"
)

func RunSmux(listenerPort int, tun *water.Interface, secretKey string, verbose bool) error {
	listener, err := net.Listen("tcp", net.JoinHostPort("", listenerPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	if verbose {
		log.Printf("Server listening on port %d\n", listenerPort)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept any connection: %v", err)
			continue
		}

		go handleSmux(conn, tun, secretKey, verbose)
	}
}

func handleSmux(conn net.Conn, tun *water.Interface, secretKey string, verbose bool) {
	defer conn.Close()

	buf := make([]byte, len(secretKey))
	_, err := conn.Read(buf)
	if err != nil {
		if verbose {
			log.Printf("Reading authentication key failed: %v", err)
		}
		return
	}

	if string(buf) != secretKey {
		if verbose {
			log.Printf("Invalid secret key")
		}
		return
	}

	smuxConfig := smux.DefaultConfig()
	session, err := smux.Server(conn, smuxConfig)
	if err != nil {
		log.Printf("Creating smux server session failed: %v", err)
		return
	}
	defer session.Close()

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Printf("Accepting stream failed: %v", err)
			return
		}
		go handleStream(stream, tun, verbose)
	}
}

func handleStream(stream net.Conn, tun *water.Interface, verbose bool) {
	defer stream.Close()

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromClient(stream, tun, clientToTun, verbose)
	go toTun(tun, clientToTun, verbose)
	go fromTun(tun, tunToClient, verbose)
	go toClient(stream, tunToClient, verbose)

	select {}
}

func fromClient(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			if verbose {
				log.Printf("Couldn't read packet length from client: %v", err)
			}
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := data(conn, buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read data from client: %v", err)
			}
			return
		}

		clientToTun <- buff
	}
}

func toTun(tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for buff := range clientToTun {
		if _, err := tun.Write(buff); err != nil {
			if verbose {
				log.Printf("Couldn't write to TUN device: %v", err)
			}
		}
	}
}

func fromTun(tun *water.Interface, tunToClient chan []byte, verbose bool) {
	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read from TUN device: %v", err)
			}
			continue
		}

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		tunToClient <- packet
	}
}

func toClient(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			if verbose {
				log.Printf("Couldn't write to client: %v", err)
			}
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
