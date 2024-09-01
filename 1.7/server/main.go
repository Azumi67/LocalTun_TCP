package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Azumi67/LocalTun_TCP/1.7/server"
	"github.com/Azumi67/LocalTun_TCP/1.7/server/utils"
)

func main() {
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("server-private", "2001:db8::1", "Server TUN IP address")
	clientIP := flag.String("client-private", "2001:db8::2", "Client TUN IP address")
	subnetMask := flag.String("subnet", "", "Subnet mask (e.g., 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	publicKeyPath := flag.String("pub-key", "public_key.pem", "public key filepath")
	mtu := flag.Int("mtu", 1500, "MTU for TUN device")
	sockEnabled := flag.Bool("sock", false, "Enable sock buffer")
	sockBuffSize := flag.Int("sockbuff", 0, "buffer size")
	keepaliveInterval := flag.Duration("keepalive", 0, "TCP keepalive interval (e.g., 10s, 1m)")
	tcpNoDelay := flag.Bool("tcpnodelay", false, "Enable TCPnodelay option")
	logEnabled := flag.Bool("log", false, "Enable logging to /etc/server.log")
	workerFlag := flag.String("worker", "0", "number of workers (default, 0 for disabled)")

	flag.Parse()

	server.RunServerSide(*serverPort, *tunIP, *clientIP, *subnetMask, *tunName, *publicKeyPath, *mtu, *sockEnabled, *sockBuffSize, *keepaliveInterval, *tcpNoDelay, *logEnabled, *workerFlag)
}
