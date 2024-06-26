package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_server"
	"github.com/Azumi67/LocalTun_TCP/monitor_server"
	"github.com/Azumi67/LocalTun_TCP/nonce_server"
	"github.com/Azumi67/LocalTun_TCP/server_smux"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/server"
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
	pingInterval := flag.Int("ping-interval", 10, "Ping interval in seconds")
	serviceName := flag.String("service-name", "azumilocal", "name of the service to restart upon failure")

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

	go monitorserver.Monitor(*clientTunIP, *pingInterval, *serviceName)
	go nonceserver.NonceDel()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warnf("Couldn't accept any connection: %v", err)
			continue
		}

		tcpnodelayserver.noDelay(conn, *tcpNoDelay)

		go handleClient(conn, tun, *secretKey, *verbose, *useSmux, *enableHeartbeat, *heartbeatInterval)
	}
}

func handleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int) {
	defer conn.Close()

	nonce, err := nonceserver.UniqueNonce()
	if err != nil {
		log.Warnf("Generating nonce failed: %v", err)
		return
	}

	if _, err := conn.Write([]byte(nonce)); err != nil {
		log.Warnf("Sending unique nonce failed: %v", err)
		return
	}

	responseBuf := make([]byte, 32)
	if _, err := conn.Read(responseBuf); err != nil {
		log.Warnf("Reading the response wasn't possible: %v", err)
		return
	}

	hashkey := sha256.New()
	hashkey.Write([]byte(nonce))
	hashkey.Write([]byte(secretKey))
	expectThisHash := hashkey.Sum(nil)

	if !nonceserver.CalHashes(responseBuf, expectThisHash) {
		log.Warnf("Wrong authentication response")
		return
	}

	if enableHeartbeat {
		go heartbeatserver.trueHeartbeat(conn, heartbeatInterval)
	}

	if useSmux {
		serversmux.HandleSmux(conn, tun, verbose)
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go tunserver.FromClient(conn, tun, clientToTun, verbose)
	go tunserver.ToTun(tun, clientToTun, verbose)
	go tunserver.FromTun(tun, tunToClient, verbose)
	go tunserver.ToClient(conn, tunToClient, verbose)

	select {}
}
//buhbye
