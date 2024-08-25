package network

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/rand"
	"encoding/binary"
	"log"
	"net"

	"github.com/songgao/water"
)

func AuthenticateClient(conn net.Conn, privateKey *rsa.PrivateKey) error {
	hashed := sha256.New()

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return err
	}

	_, err = conn.Write(signature)
	if err != nil {
		return err
	}

	return nil
}

func HandleConnection(conn net.Conn, tun *water.Interface, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		_, err := conn.Read(pcktLength)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read packet length: %v", err)
			}
			continue
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)
		n, err := conn.Read(buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read data from the server: %v", err)
			}
			continue
		}

		if n != int(length) {
			if verbose {
				log.Printf("Packet mismatch: expected %d, gotten %d", length, n)
			}
			continue
		}

		_, err = tun.Write(buff[:n])
		if err != nil {
			if verbose {
				log.Printf("Couldn't write to TUN device: %v", err)
			}
			continue
		}

		if verbose {
			log.Printf("Written %d bytes to TUN device", n)
		}
	}
}

func ForwardTraffic(tun *water.Interface, conn net.Conn, verbose bool) {
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

		_, err = conn.Write(packet)
		if err != nil {
			if verbose {
				log.Printf("Couldn't write to server: %v", err)
			}
			continue
		}

		if verbose {
			log.Printf("Forwarded %d bytes to the server", n)
		}
	}
}
