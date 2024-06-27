// main.go
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_server"
	"github.com/Azumi67/LocalTun_TCP/monitor_server"
	"github.com/Azumi67/LocalTun_TCP/nonce_server"
	"github.com/Azumi67/LocalTun_TCP/server_smux"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/tun_server"
	"github.com/Azumi67/LocalTun_TCP/worker_server"
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
	workers := flag.Int("worker", runtime.NumCPU(), "number of workers or based on the number of CPU cores")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	tun, err := tunserver.TunUp(*tunIP, *clientTunIP, *subnetMask, *mtu, *tunName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer tun.Close()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", *serverPort, err)
	}
	defer ln.Close()

	log.Infof("Server listening on port %d", *serverPort)

	go monitor_server.Monitor(*clientTunIP, *pingInterval, *serviceName)
	go nonce_server.NonceDel()

	var wg sync.WaitGroup
	workerOnechan := make(chan net.Conn, *workers)

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker_server.Worker(tun, workerOnechan, *secretKey, *verbose, *useSmux, *enableHeartbeat, *heartbeatInterval, &wg)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warnf("Couldn't accept any connection: %v", err)
			continue
		}
		if *tcpNoDelay {
			tcp_no_delay_server.noDelay(conn)
		}
		workerOnechan <- conn
	}

	close(workerOnechan)
	wg.Wait()
}
