package network

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/songgao/water"
)

func Handle(conn net.Conn, tun *water.Interface, publicKey *rsa.PublicKey, verbose bool) {
	defer conn.Close()

	hashed := sha256.New()

	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("Couldn't read authentication key: %v", err)
		return
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed.Sum(nil), buf); err != nil {
		log.Printf("invalid sign")
		return
	}

	log.Println("Client authenticated")

	go func() {
		for {
			pcktLength := make([]byte, 2)
			_, err := conn.Read(pcktLength)
			if err != nil {
				continue
			}

			length := binary.BigEndian.Uint16(pcktLength)
			buff := make([]byte, length)
			n, err := conn.Read(buff)
			if err != nil {
				log.Printf("Couldn't read data from client: %v", err)
				continue
			}

			if n != int(length) {
				log.Printf("Packet mismatch: i've expected %d, but gotten %d", length, n)
				continue
			}

			_, err = tun.Write(buff[:n])
			if err != nil {
				log.Printf("Couldn't write to TUN device: %v", err)
				continue
			}

			if verbose {
				log.Printf("Written %d bytes to TUN device\n", n)
			}
		}
	}()

	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			log.Printf("Couldn't read from TUN device: %v", err)
			continue
		}

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		_, err = conn.Write(packet)
		if err != nil {
			log.Printf("Couldn't write to client: %v", err)
			continue
		}

		if verbose {
			log.Printf("Forwarded %d bytes to the client\n", n)
		}
	}
}
