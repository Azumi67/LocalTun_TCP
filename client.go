package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
        "github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/client"
	"github.com/Azumi67/LocalTun_TCP/smux_client"
)

var log = logrus.New()

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server Public IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	clientTunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g: 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")
	tcpNoDelay := flag.Bool("tcp-nodelay", false, "Enable TCP_NODELAY")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	tun, err := tunclient.tunUp(*tunName, *clientTunIP, *serverTunIP, *subnetMask, *mtu)
	if err != nil {
		log.Fatalf("Couldn't setup TUN device: %v", err)
	}
	defer tun.Close()

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *kharejAddr, *serverPort))
		if err != nil {
			log.Warnf("Couldn't connect to server: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		if err := tcpnodelayclient.NoDelay(conn, *tcpNoDelay); err != nil {
			log.Warnf("Couldn't set TCP_NODELAY: %v", err)
		}

		log.Info("Connected to server")

		stop := make(chan struct{})
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			if *useSmux {
				smuxclient.HandleSmux(conn, tun, *verbose, stop)
			} else {
				tunclient.HandleServer(conn, tun, *secretKey, *verbose, *useSmux, stop)
			}
		}()

		for {
			time.Sleep(10 * time.Second)
			if !ping(*serverTunIP) {
				log.Warn("Server is unreachable, Azumi is reconnecting..")
				close(stop)
				conn.Close()
				wg.Wait()
				break
			}
		}
	}
}

func ping(address string) bool {
	cmd := exec.Command("ping", "-c1", "-w1", address)
	err := cmd.Run()
	return err == nil
}
