package main

import (
	"flag"
	"log"

	"github.com/Azumi67/LocalTun_TCP/client"
)

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "10.0.0.2", "TUN IP address")
	kharejTunIP := flag.String("server-private", "10.0.0.1", "Server TUN IP address")
	subnetMask := flag.String("subnet", "24", "Subnet mask (e.g., 24 or 64)")
	tunName := flag.String("device", "tun0", "TUN device name")
	secretKey := flag.String("key", "", "Secret key for authentication")
	mtu := flag.Int("mtu", 1500, "MTU for TUN device")

	flag.Parse()

	config := tunclient.ClientConfig{
		ServerAddr:  *kharejAddr,
		ServerPort:  *serverPort,
		TunIP:       *tunIP,
		KharejTunIP: *kharejTunIP,
		SubnetMask:  *subnetMask,
		TunName:     *tunName,
		SecretKey:   *secretKey,
		MTU:         *mtu,
	}

	log.Printf("Client with config Started: %+v\n", config)

	tunclient.clientSide(config)
}
