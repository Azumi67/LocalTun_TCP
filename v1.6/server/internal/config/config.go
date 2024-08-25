package config

import (
	"flag"
)

type ServerConfig struct {
	ServerPort        int
	TunIP             string
	ClientTunIP       string
	SubnetMask        string
	TunName           string
	PublicKeyPath     string
	MTU               int
	EnableSmux        bool
	EnableTCPNoDelay  bool
	KeepAliveInterval int
	Verbose           bool
}

func LoadServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	flag.IntVar(&cfg.ServerPort, "server-port", 808, "Server port")
	flag.StringVar(&cfg.TunIP, "server-private", "2001:db8::1", "Server TUN IP address")
	flag.StringVar(&cfg.ClientTunIP, "client-private", "2001:db8::2", "Client TUN IP address")
	flag.StringVar(&cfg.SubnetMask, "subnet", "64", "Subnet mask (24 or 64)")
	flag.StringVar(&cfg.TunName, "device", "tun2", "TUN device name")
	flag.StringVar(&cfg.PublicKeyPath, "pub-key", "public_key.pem", "Path to the public key file")
	flag.IntVar(&cfg.MTU, "mtu", 1500, "MTU for TUN device")
	flag.BoolVar(&cfg.EnableSmux, "smux", false, "Enable smux for multiplexing")
	flag.BoolVar(&cfg.EnableTCPNoDelay, "tcpnodelay", false, "Enabling tcpnodelay option")
	flag.IntVar(&cfg.KeepAliveInterval, "keepalive", 0, "Enable TCP keepalive with custom interval in seconds")
	flag.BoolVar(&cfg.Verbose, "v", false, "Enabling verbose logging")

	flag.Parse()

	return cfg
}
