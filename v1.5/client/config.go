package client

import (
	"flag"
	"runtime"
)

type Config struct {
	ServerAddr string
	ServerPort int
	TunIP      string
	KharejTunIP string
	SubnetMask string
	TunName    string
	SecretKey  string
	MTU        int
	Verbose    bool
	TCPNoDelay bool
	NumWorkers int
}

func ParseConfig() *Config {
	serverAddr := flag.String("server-addr", "SERVER_IP", "Server IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "2001:db8::2", "TUN IP address")
	kharejTunIP := flag.String("server-private", "2001:db8::1", "Server TUN IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g., 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1430, "MTU for TUN device")
	verbose := flag.Bool("v", false, "enable verbose logging")
	tcpNoDelay := flag.Bool("tcpnodelay", false, "disable Nagle's algorithm")
	numWorkers := flag.Int("worker", runtime.NumCPU(), "Number of worker's goroutines")

	flag.Parse()

	if *numWorkers == 0 {
		*numWorkers = runtime.NumCPU()
	}

	return &Config{
		ServerAddr: *serverAddr,
		ServerPort: *serverPort,
		TunIP:      *tunIP,
		KharejTunIP: *kharejTunIP,
		SubnetMask: *subnetMask,
		TunName:    *tunName,
		SecretKey:  *secretKey,
		MTU:        *mtu,
		Verbose:    *verbose,
		TCPNoDelay: *tcpNoDelay,
		NumWorkers: *numWorkers,
	}
}
