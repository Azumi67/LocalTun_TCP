package server

import (
	"flag"
	"runtime"
)

type Config struct {
	ServerPort int
	TunIP      string
	IranIP     string
	SubnetMask string
	TunName    string
	SecretKey  string
	MTU        int
	Verbose    bool
	TCPNoDelay bool
	NumWorkers int
}

func ParseConfig() *Config {
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("server-private", "2001:db8::1", "Kharej IP address")
	iranIP := flag.String("client-private", "2001:db8::2", "Client IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1500, "MTU for TUN device")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	tcpNoDelay := flag.Bool("tcpnodelay", false, "Disable Nagle's algorithm")
	numWorkers := flag.Int("worker", runtime.NumCPU(), "Number of worker goroutines")

	flag.Parse()

	if *numWorkers == 0 {
		*numWorkers = runtime.NumCPU()
	}

	return &Config{
		ServerPort: *serverPort,
		TunIP:      *tunIP,
		IranIP:     *iranIP,
		SubnetMask: *subnetMask,
		TunName:    *tunName,
		SecretKey:  *secretKey,
		MTU:        *mtu,
		Verbose:    *verbose,
		TCPNoDelay: *tcpNoDelay,
		NumWorkers: *numWorkers,
	}
}
