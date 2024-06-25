package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/server"
	"github.com/Azumi67/LocalTun_TCP/server_smux"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_server"
)

var log = logrus.New()

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
	clientTunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g: 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")
	enableHeartbeat := flag.Bool("heartbeat", false, "Enable heartbeat")
	heartbeatInterval := flag.Int("heartbeat-interval", 30, "Heartbeat interval in seconds")
	tcpNoDelay := flag.Bool("tcp-nodelay", false, "Enable TCP_NODELAY")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = *tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunserver.TunUp(tun, *tunIP, *clientTunIP, *subnetMask, *mtu, *tunName)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", *serverPort, err)
	}
	defer ln.Close()

	log.Infof("Server listening on port %d", *serverPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Couldn't accept any connection: %v", err)
			continue
		}

		if *tcpNoDelay {
			tcpnodelayserver.noDelay(conn)
		}

		go handleClient(conn, tun, *secretKey, *verbose, *useSmux, *enableHeartbeat, *heartbeatInterval)
	}
}

func handleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int) {
	defer conn.Close()

	authBuf := make([]byte, len(secretKey))
	if _, err := conn.Read(authBuf); err != nil {
		log.Warnf("Failed to read authentication key: %v", err)
		return
	}

	if string(authBuf) != secretKey {
		log.Warnf("Wrong authentication key")
		return
	}

	if enableHeartbeat {
		heartbeatserver.trueHeartbeat(conn, heartbeatInterval)
	}

	if useSmux {
		serversmux.HandleSmux(conn, tun, verbose)
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromClient(conn, tun, clientToTun, verbose)
	go toTun(tun, clientToTun, verbose)
	go fromTun(tun, tunToClient, verbose)
	go toClient(conn, tunToClient, verbose)

	select {}
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

func toClient(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			log.Warnf("Couldn't write to client: %v", err)
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
