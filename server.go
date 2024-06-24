package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/server_smux"
	"github.com/Azumi67/LocalTun_TCP/server"
)

var log = logrus.New()

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
	clientTunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g: 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable SMUX multiplexing")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	tun, err := tunserver.Run(*serverTunIP, *clientTunIP, *subnetMask, *tunName, *mtu)
	if err != nil {
		log.Fatalf("failed to set up TUN: %v", err)
	}
	defer tun.Close()

	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
		if err != nil {
			log.Errorf("failed to listen on port %d: %v", *serverPort, err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer ln.Close()

		log.Infof("Server listening on port %d", *serverPort)

		conn, err := ln.Accept()
		if err != nil {
			log.Warnf("failed to accept any connection: %v", err)
			ln.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info("Client connected")

		if *useSmux {
			smuxserver.HandleSmux(conn, tun, *verbose)
		} else {
			tcpnodelayserver.noDelay(conn)
			handleClient(conn, tun, *secretKey, *verbose)
		}
	}
}

func handleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool) {
	authBuf := make([]byte, len(secretKey))
	if _, err := conn.Read(authBuf); err != nil {
		log.Warnf("failed to read authentication key: %v", err)
		return
	}

	if string(authBuf) != secretKey {
		log.Warnf("Wrong authentication key")
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
//byebye
