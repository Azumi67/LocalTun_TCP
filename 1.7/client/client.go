package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/1.7/client/utils"
)

func RunClientSide(serverAddr string, serverPort int, tunIP, serverTunIP, subnetMask, tunName, privateKeyPath string, mtu int, sockEnabled bool, sockBuffSize int, keepaliveInterval time.Duration, tcpNoDelay, logEnabled bool, workerFlag string) {
	if subnetMask == "" {
		subnetMask = utils.DefaultSubnet(tunIP)
	}

	if logEnabled {
		logFile, err := utils.LogFile("/etc/client.log")
		if err != nil {
			log.Fatalf("Couldn't open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	privateKey, err := utils.LoadPrivateKey(privateKeyPath)
	if err != nil {
		log.Fatalf("Couldn't load private key: %v", err)
	}

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	log.Printf("TUN device created: %s\n", tun.Name())

	if err := utils.Cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := utils.Cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := utils.Cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if utils.IPv6(serverTunIP) {
		if err := utils.Cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", serverTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := utils.Cmd("ip", "route", "add", fmt.Sprintf("%s/32", serverTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}

	if err := utils.Cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward"); err != nil {
		log.Fatalf("Enabling IPv4 forwarding failed: %v", err)
	}
	if err := utils.Cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding"); err != nil {
		log.Fatalf("Enabling IPv6 forwarding failed: %v", err)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("[%s]:%d", serverAddr, serverPort))
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	if keepaliveInterval > 0 {
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(keepaliveInterval)
		}
	}

	if sockEnabled && sockBuffSize > 0 {
		err := conn.(*net.TCPConn).SetReadBuffer(sockBuffSize)
		if err != nil {
			log.Printf("set read buffer size failed: %v", err)
		}
		err = conn.(*net.TCPConn).SetWriteBuffer(sockBuffSize)
		if err != nil {
			log.Printf("set write buffer size failed: %v", err)
		}
	}

	if tcpNoDelay {
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)
		}
	}

	hashed := sha256.New()

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		log.Fatalf("signing data wasn't possible: %v", err)
	}

	_, err = conn.Write(signature)
	if err != nil {
		log.Fatalf("Couldn't send auth key: %v", err)
	}

	workerCount := utils.WorkerCount(workerFlag)

	if workerCount > 0 {

		taskChan := make(chan []byte, workerCount)
		var workerWg sync.WaitGroup

		for i := 0; i < workerCount; i++ {
			workerWg.Add(1)
			go worker(tun, taskChan, &workerWg)
		}

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
					log.Printf("Couldn't read data from the server: %v", err)
					continue
				}

				if n != int(length) {
					log.Printf("Packet length mismatch: expected %d, but gotten %d", length, n)
					continue
				}

				taskChan <- buff[:n]
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
				log.Printf("Couldn't write to server: %v", err)
				continue
			}

			log.Printf("Forwarded %d bytes to the server\n", n)
		}

		close(taskChan)
		workerWg.Wait()
	} else {

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
				log.Printf("Couldn't read data from the server: %v", err)
				continue
			}

			if n != int(length) {
				log.Printf("Packet length mismatch: expected %d, but gotten %d", length, n)
				continue
			}

			_, err = tun.Write(buff[:n])
			if err != nil {
				log.Printf("Couldn't write to TUN device: %v", err)
				continue
			}

			log.Printf("Written %d bytes to TUN device\n", n)
		}
	}
}

func worker(tun *water.Interface, taskChan <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		_, err := tun.Write(task)
		if err != nil {
			log.Printf("Couldn't write to TUN device: %v", err)
			continue
		}
		log.Printf("Worker task: written %d bytes to TUN device\n", len(task))
	}
}
