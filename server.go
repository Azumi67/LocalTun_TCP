package main

import (
	"flag"
	"log"
	"time"

	"github.com/Azumi67/LocalTun_TCP/server"
)

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("server-private", "2001:db8::1", "Kharej IP address")
	iranIP := flag.String("client-private", "2001:db8::2", "Client IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")

	flag.Parse()

	tunserver.Run(*serverPort, *tunIP, *iranIP, *subnetMask, *tunName, *mtu, *secretKey, *verbose)
}
