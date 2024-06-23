package tunclient

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/songgao/water"
)


func Run(tun *water.Interface, serverAddr string, serverPort int, secretKey string, verbose bool) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("[%s]:%d", serverAddr, serverPort))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	if verbose {
		log.Println("Connected to server")
	}

	_, err = conn.Write([]byte(secretKey))
	if err != nil {
		return fmt.Errorf("failed to send authentication key: %v", err)
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromServer(conn, tun, clientToTun, verbose)
	go toTun(tun, clientToTun, verbose)
	go fromTun(tun, tunToClient, verbose)
	go toServer(conn, tunToClient, verbose)

	select {}
}

func fromServer(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			if verbose {
				log.Printf("Couldn't read packet length from server: %v", err)
			}
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := data(conn, buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read data from server: %v", err)
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

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			if verbose {
				log.Printf("Couldn't write to server: %v", err)
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

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
