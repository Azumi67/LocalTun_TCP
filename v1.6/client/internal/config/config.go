package config

import (
	"flag"
)

type ClientConfig struct {
	ServerAddr        string
	ServerPort        int
	TunIP             string
	ServerTunIP       string
	SubnetMask        string
	TunName           string
	PrivateKeyPath    string
	MTU               int
	EnableSmux        bool
	EnableTCPNoDelay  bool
	KeepAliveInterval int
	Verbose           bool
}

func LoadClientConfig() *ClientConfig {
	cfg := &ClientConfig{}
	flag.StringVar(&cfg.ServerAddr, "server-addr", "SERVER_IP", "Server IP address")
	flag.IntVar(&cfg.ServerPort, "server-port", 800, "Server port")
	flag.StringVar(&cfg.TunIP, "client-private", "2001:db8::2", "Client TUN IP address")
	flag.StringVar(&cfg.ServerTunIP, "server-private", "2001:db8::1", "Server TUN IP address")
	flag.StringVar(&cfg.SubnetMask, "subnet", "64", "Subnet mask (24 or 64)")
	flag.StringVar(&cfg.TunName, "device", "tun2", "TUN device name")
	flag.StringVar(&cfg.PrivateKeyPath, "priv-key", "private_key.pem", "private key path file")
	flag.IntVar(&cfg.MTU, "mtu", 1430, "MTU for TUN device")
	flag.BoolVar(&cfg.EnableSmux, "smux", false, "Enable smux for multiplexing")
	flag.BoolVar(&cfg.EnableTCPNoDelay, "tcpnodelay", false, "Enabling tcpnodelay option")
	flag.IntVar(&cfg.KeepAliveInterval, "keepalive", 0, "Enable TCP keepalive with custom interval in seconds")
	flag.BoolVar(&cfg.Verbose, "v", false, "Enabling verbose logging")

	flag.Parse()

	return cfg
}
