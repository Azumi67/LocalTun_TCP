package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/client"
)

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "2001:db8::2", "TUN IP address")
	kharejTunIP := flag.String("server-private", "2001:db8::1", "Server TUN IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g., 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1400, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")

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

	tunclient.tunUP(tun, *tunIP, *kharejTunIP, *subnetMask, *mtu, *tunName)

	for {
		err = clientpackage.clientSide(tun, *kharejAddr, *serverPort, *secretKey, *verbose)
		if err != nil {
			log.Printf("error in client loop: %v. Retrying in 5 seconds..\n", err)
			time.Sleep(5 * time.Second)
		}
	}
}
