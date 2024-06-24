package tunserver

import (
	"fmt"
	"os/exec"
	"net"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

var log = logrus.New()

func Run(serverTunIP, clientTunIP, subnetMask, tunName string, mtu int) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = tunName
	tun, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("couldn't create TUN device: %v", err)
	}

	if err := tunUp(tun, serverTunIP, clientTunIP, subnetMask, mtu); err != nil {
		return nil, err
	}

	return tun, nil
}

func tunUp(tun *water.Interface, serverTunIP, clientTunIP, subnetMask string, mtu int) error {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		return fmt.Errorf("couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		return fmt.Errorf("couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", serverTunIP, subnetMask), "dev", tun.Name()); err != nil {
		return fmt.Errorf("couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(clientTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", clientTunIP), "dev", tun.Name()); err != nil {
			return fmt.Errorf("adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", clientTunIP), "dev", tun.Name()); err != nil {
			return fmt.Errorf("adding route for private IPv4 failed: %v", err)
		}
	}

	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward"); err != nil {
		return fmt.Errorf("enabling IPv4 forwarding failed: %v", err)
	}
	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding"); err != nil {
		return fmt.Errorf("enabling IPv6 forwarding failed: %v", err)
	}

	return nil
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func iPv6(s string) bool {
	return net.ParseIP(s).To4() == nil
}
