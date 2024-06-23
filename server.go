package main

import (
	"flag"
	"log"
        "github.com/Azumi67/LocalTun_TCP/server_smux"
	"github.com/Azumi67/LocalTun_TCP/server"
	"github.com/songgao/water"
)


func main() {
	serverPort := flag.Int("server-port", 8000, "Server port")
	tunIP := flag.String("server-private", "2001:db8::1", "Server TUN IP address")
	clientIP := flag.String("client-private", "2001:db8::2", "Client TUN IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")

	flag.Parse()

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = *tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunUp(tun, *tunIP, *clientIP, *subnetMask, *mtu, *tunName)

	if *useSmux {
		err := smuxserver.RunSmux(*serverPort, tun, *secretKey, *verbose)
		if err != nil {
			log.Fatalf("Error running SMUX server: %v", err)
		}
	} else {
		err := tunserver.Run(*serverPort, *tunIP, *clientIP, *subnetMask, *tunName, *secretKey, *mtu, *verbose)
		if err != nil {
			log.Fatalf("Error running TUN server: %v", err)
		}
	}
}

func tunUp(tun *water.Interface, tunIP, clientIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(clientIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", clientIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", clientIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}

	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward"); err != nil {
		log.Fatalf("Enabling IPv4 forwarding failed: %v", err)
	}
	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding"); err != nil {
		log.Fatalf("Enabling IPv6 forwarding failed: %v", err)
	}
}
