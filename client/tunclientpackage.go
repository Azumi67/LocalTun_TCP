package tunclient

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

func tunUP(tun *water.Interface, tunIP, kharejTunIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(kharejTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", kharejTunIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", kharejTunIP), "dev", tun.Name()); err != nil {
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

func fromServer(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			if verbose {
				log.Printf("Couldn't read packet length from server: %v", err)
			}
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := data(conn, buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read data from server: %v", err)
			}
			return
		}

		clientToTun <- buff
	}
}

func toTun(tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for buff := range clientToTun {
		if _, err := tun.Write(buff); err != nil {
			if verbose {
				log.Printf("Couldn't write to TUN device: %v", err)
			}
		}
	}
}

func fromTun(tun *water.Interface, tunToClient chan []byte, verbose bool) {
	for {
		buff := make([]byte, 1500)
		n, err := tun.Read(buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read from TUN device: %v", err)
			}
			continue
		}

		packet := make([]byte, 2+n)
		binary.BigEndian.PutUint16(packet[:2], uint16(n))
		copy(packet[2:], buff[:n])

		tunToClient <- packet
	}
}

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			if verbose {
				log.Printf("Couldn't write to server: %v", err)
			}
		}
	}
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func data(conn net.Conn, buff []byte) (int, error) {
	totalRead := 0
	for totalRead < len(buff) {
		n, err := conn.Read(buff[totalRead:])
		if err != nil {
			return totalRead, err
		}
		totalRead += n
	}
	return totalRead, nil
}

func iPv6(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil
}
