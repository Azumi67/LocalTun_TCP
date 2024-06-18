package tunclient

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

type ClientConfig struct {
	ServerAddr  string
	ServerPort  int
	TunIP       string
	KharejTunIP string
	SubnetMask  string
	TunName     string
	SecretKey   string
	MTU         int
}

func clientSide(config ClientConfig) {
	tun, err := tuncreate(config.TunName)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunup(tun, config.TunIP, config.SubnetMask, config.MTU)

	if iPv6(config.KharejTunIP) {
		if err := route(tun.Name(), config.KharejTunIP, "128"); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := route(tun.Name(), config.KharejTunIP, "32"); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}

	enableForward()

	conn, err := net.Dial("tcp", fmt.Sprintf("[%s]:%d", config.ServerAddr, config.ServerPort))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	_, err = conn.Write([]byte(config.SecretKey))
	if err != nil {
		log.Fatalf("Failed to send authentication key: %v", err)
	}

	go pckttotun(tun, conn)
	pckttokharej(tun, conn)
}

func tuncreate(name string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name
	return water.New(config)
}

func tunup(tun *water.Interface, tunIP, subnetMask string, mtu int) {
	log.Printf("TUN device created: %s\n", tun.Name())

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	log.Printf("TUN device is up: %s\n", tun.Name())

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	log.Printf("Private IP address %s/%s assigned to TUN device: %s\n", tunIP, subnetMask, tun.Name())
}

func route(deviceName, ip, mask string) error {
	if iPv6(ip) {
		return cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/%s", ip, mask), "dev", deviceName)
	} else {
		return cmd("ip", "route", "add", fmt.Sprintf("%s/%s", ip, mask), "dev", deviceName)
	}
}

func enableForward() {
	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward"); err != nil {
		log.Fatalf("Enabling IPv4 forwarding failed: %v", err)
	}
	if err := cmd("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding"); err != nil {
		log.Fatalf("Enabling IPv6 forwarding failed: %v", err)
	}
}

func pckttotun(tun *water.Interface, conn net.Conn) {
	for {
		pcktLength := make([]byte, 2)
		_, err := conn.Read(pcktLength)
		if err != nil {
			log.Fatalf("Couldn't read packet from the server: %v", err)
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatalf("Couldn't read data from the server: %v", err)
		}

		log.Printf("Reading %d bytes from the server\n", n)

		_, err = tun.Write(buff[:n])
		if err != nil {
			log.Fatalf("Couldn't write to TUN device: %v", err)
		}

		log.Printf("Written %d bytes to TUN device\n", n)
	}
}

func pckttokharej(tun *water.Interface, conn net.Conn) {
	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			log.Fatalf("Couldn't read from TUN device: %v", err)
		}

		log.Printf("Reading %d bytes from TUN device\n", n)

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		_, err = conn.Write(packet)
		if err != nil {
			log.Fatalf("Couldn't write to server: %v", err)
		}

		log.Printf("Forwarded %d bytes to the server\n", n)
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
