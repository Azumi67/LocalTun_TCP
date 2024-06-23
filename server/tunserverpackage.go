package tunserver

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/songgao/water"
)

func Run(serverPort int, tunIP, clientIP, subnetMask, tunName, secretKey string, mtu int, verbose bool) error {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = tunName
	tun, err := water.New(config)
	if err != nil {
		return fmt.Errorf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunUp(tun, tunIP, clientIP, subnetMask, mtu, tunName)

	for {
		err = serverSide(tun, serverPort, secretKey, verbose)
		if err != nil {
			log.Printf("Server loop error: %v. Retrying in 3 seconds..\n", err)
			time.Sleep(3 * time.Second)
		}
	}
}

func tunUp(tun *water.Interface, tunIP, clientIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(clientIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", clientIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", clientIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}

	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward"); err != nil {
		log.Fatalf("Enabling IPv4 forwarding failed: %v", err)
	}
	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding"); err != nil {
		log.Fatalf("Enabling IPv6 forwarding failed: %v", err)
	}
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func iPv6(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil
}
