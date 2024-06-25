package tunclient

import (
	"fmt"
	"os/exec"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"net"
)

var log = logrus.New()

func TunUp(tun *water.Interface, tunIP, serverTunIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if isIPv6(serverTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", serverTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", serverTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}
}

func cmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func isIPv6(ip string) bool {
	return net.ParseIP(ip).To4() == nil
}
