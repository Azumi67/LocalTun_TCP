package tunserver

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

type kharejConfig struct {
	ServerPort int
	TunIP      string
	IranIP     string
	SubnetMask string
	TunName    string
	SecretKey  string
	MTU        int
}

func serverSide(config kharejConfig) {
	tun, err := tuncreate(config.TunName)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunup(tun, config.TunIP, config.SubnetMask, config.MTU)

	if iPv6(config.IranIP) {
		if err := route(tun.Name(), config.IranIP, "128"); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := route(tun.Name(), config.IranIP, "32"); err != nil {
			log.Fatalf("Adding route for private IPv4 failed: %v", err)
		}
	}

	enableForward()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", config.ServerPort, err)
	}
	defer listener.Close()

	log.Printf("Server is listening on port %d\n", config.ServerPort)

	iranFunctions(listener, tun, config.SecretKey)
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

func route(devName, ip, mask string) error {
	if iPv6(ip) {
		return cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/%s", ip, mask), "dev", devName)
	} else {
		return cmd("ip", "route", "add", fmt.Sprintf("%s/%s", ip, mask), "dev", devName)
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

func iranFunctions(listener net.Listener, tun *water.Interface, secretKey string) {
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Couldn't accept TCP connection: %v", err)
	}
	defer conn.Close()

	log.Println("Client connected")

	auth := make([]byte, len(secretKey))
	_, err = conn.Read(auth)
	if err != nil {
		log.Fatalf("Reading secret key from client failed: %v", err)
	}

	if string(auth) != secretKey {
		log.Fatalf("Invalid secret key from client")
	}

	log.Println("Client authenticated successfully")

	go pckttoiran(tun, conn)
	pckttotun(tun, conn)
}

func pckttoiran(tun *water.Interface, conn net.Conn) {
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
			log.Fatalf("Couldn't write to client: %v", err)
		}

		log.Printf("Forwarded %d bytes to client\n", n)
	}
}

func pckttotun(tun *water.Interface, conn net.Conn) {
	for {
		pcktLength := make([]byte, 2)
		_, err := conn.Read(pcktLength)
		if err != nil {
			log.Fatalf("Reading packet from client failed: %v", err)
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatalf("Couldn't read from client: %v", err)
		}

		log.Printf("Reading %d bytes from client\n", n)

		_, err = tun.Write(buff[:n])
		if err != nil {
			log.Fatalf("Couldn't write to TUN device: %v", err)
		}

		log.Printf("Written %d bytes to TUN device\n", n)
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
