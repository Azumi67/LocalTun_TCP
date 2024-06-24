package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_server"
	"github.com/Azumi67/LocalTun_TCP/server_smux"
)

var log = logrus.New()

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server private IP address")
	clientTunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g: 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")
	tcpNoDelay := flag.Bool("tcp-nodelay", false, "Enable TCPNODELAY")

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

	tunUp(tun, *serverTunIP, *clientTunIP, *subnetMask, *mtu, *tunName)

	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
		if err != nil {
			log.Errorf("failed to listen on port %d: %v", err)
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
			serversmux.HandleClient(conn, tun, *secretKey, *verbose)
		} else {
			tcpnodelayserver.HandleClient(conn, tun, *secretKey, *verbose, *tcpNoDelay)
		}
	}
}

func tunUp(tun *water.Interface, serverTunIP, clientTunIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", serverTunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(clientTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", clientTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", clientTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func iPv6(address string) bool {
	return net.ParseIP(address).To4() == nil
}
