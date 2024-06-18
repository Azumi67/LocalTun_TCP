package main

import (
	"flag"
	"log"

	"github.com/Azumi67/LocalTun_TCP/server"
)

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("server-private", "10.0.0.1", "Kharej IP address")
	iranIP := flag.String("client-private", "10.0.0.2", "Client IP address")
	subnetMask := flag.String("subnet", "24", "Subnet mask (24 or 64)")
	tunName := flag.String("device", "tun0", "TUN device name")
	secretKey := flag.String("key", "", "Secret key for authentication")
	mtu := flag.Int("mtu", 1500, "MTU for TUN device")

	flag.Parse()

	config := tunserver.kharejConfig{
		ServerPort: *serverPort,
		TunIP:      *tunIP,
		IranIP:     *iranIP,
		SubnetMask: *subnetMask,
		TunName:    *tunName,
		SecretKey:  *secretKey,
		MTU:        *mtu,
	}

	log.Printf("Server with config Started: %+v\n", config)

	tunserver.serverSide(config)
}
