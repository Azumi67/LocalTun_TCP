package client

import (
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"net"
	"github.com/Azumi67/LocalTun_TCP/v1.5/utils"
	"io"
	"encoding/binary"
)

func Worker(log *logrus.Logger, conn net.Conn, tun *water.Interface, config *Config) {
	clientNonce := utils.GenerateNonce(log)
	if clientNonce == nil {
		return
	}

	if _, err := conn.Write(clientNonce); err != nil {
		log.Fatalf("Couldn't send Client's nonce: %v", err)
	}

	serverNonce := make([]byte, 32)
	if _, err := io.ReadFull(conn, serverNonce); err != nil {
		log.Fatalf("Couldn't read Server's nonce: %v", err)
	}

	clientHash := utils.GenerateHash(clientNonce, serverNonce, config.SecretKey)
	if _, err := conn.Write(clientHash); err != nil {
		log.Fatalf("failed to send client's hashkey: %v", err)
	}

	serverHash := make([]byte, 32)
	if _, err := io.ReadFull(conn, serverHash); err != nil {
		log.Fatalf("Couldn't read server's hashkey: %v", err)
	}
	expectedServerHash := utils.GenerateHash(serverNonce, clientNonce, config.SecretKey)
	if !utils.CompareHashes(serverHash, expectedServerHash) {
		log.Warn("Wrong server's hashkey")
		return
	}

	log.Info("Server authenticated successfully")

	go func() {
		for {
			pcktLength := make([]byte, 2)
			_, err := conn.Read(pcktLength)
			if err != nil {
				if err == io.EOF {
					log.Info("Server closed the connection")
					return
				}
				log.Errorf("Couldn't read packet length from server: %v", err)
				return
			}

			length := binary.BigEndian.Uint16(pcktLength)
			buff := make([]byte, length)
			n, err := conn.Read(buff)
			if err != nil {
				log.Errorf("Couldn't read data from server: %v", err)
				return
			}

			if n != int(length) {
				//log.Warnf("Packet length mismatch: expected %d, got %d", length, n)
				continue
			}

			_, err = tun.Write(buff[:n])
			if err != nil {
				log.Errorf("Couldn't write to TUN device: %v", err)
				continue
			}

			log.Infof("Written %d bytes to TUN device", n)
		}
	}()

	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			log.Errorf("Couldn't read from TUN device: %v", err)
			continue
		}

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		_, err = conn.Write(packet)
		if err != nil {
			if err == io.EOF {
				log.Info("Server closed the connection")
				return
			}
			log.Errorf("Couldn't write to server: %v", err)
			continue
		}

		log.Infof("Forwarded %d bytes to the server", n)
	}
}
