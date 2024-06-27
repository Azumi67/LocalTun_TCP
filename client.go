package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/client"
	"github.com/Azumi67/LocalTun_TCP/smux_client"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_client"
	"github.com/Azumi67/LocalTun_TCP/monitor_client"
	"github.com/Azumi67/LocalTun_TCP/hash"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_client"
	"github.com/Azumi67/LocalTun_TCP/worker_client"
)

var log = logrus.New()

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server Public IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
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

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = *tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunclient.TunUp(tun, *tunIP, *serverTunIP, *subnetMask, *mtu, *tunName)

	go monitorclient.Monitor(*serverTunIP, *pingInterval, *serviceName)

	var wg sync.WaitGroup
	workerOnechan := make(chan net.Conn, *workers)

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go workerclient.Worker(tun, workerOnechan, *secretKey, *verbose, *useSmux, *enableHeartbeat, *heartbeatInterval, &wg)
	}

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *kharejAddr, *serverPort))
		if err != nil {
			log.Warnf("Couldn't connect to server: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if *tcpNoDelay {
			tcpnodelayclient.noDelay(conn)
		}

		defer conn.Close()

		log.Info("Connected to server")
		workerOnechan <- conn
	}

	close(workerOnechan)
	wg.Wait()
}
