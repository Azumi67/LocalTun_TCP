package network

import (
	"fmt"
	"os/exec"

	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/v1.6/client/internal/config"
)

func SetupTunDevice(cfg *config.ClientConfig) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = cfg.TunName
	tun, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("couldn't create TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		return nil, fmt.Errorf("couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", cfg.MTU)); err != nil {
		return nil, fmt.Errorf("couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", cfg.TunIP, cfg.SubnetMask), "dev", tun.Name()); err != nil {
		return nil, fmt.Errorf("couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(cfg.ServerTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", cfg.ServerTunIP), "dev", tun.Name()); err != nil {
			return nil, fmt.Errorf("adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", cfg.ServerTunIP), "dev", tun.Name()); err != nil {
			return nil, fmt.Errorf("adding route for private IPv4 failed: %v", err)
		}
	}

	return tun, nil
}
