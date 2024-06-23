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

func Run(serverPort int, tunIP, iranIP, subnetMask string, tunName string, mtu int, secretKey string, verbose bool) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunUp(tun, tunIP, iranIP, subnetMask, mtu, tunName)

	for {
		err = serverSide(tun, serverPort, secretKey, verbose)
		if err != nil {
			log.Printf("Server loop error: %v. Retrying in 3 seconds...\n", err)
			time.Sleep(3 * time.Second)
		}
	}
}

func tunUp(tun *water.Interface, tunIP, iranIP, subnetMask string, mtu int, tunName string) {
	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		log.Fatalf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		log.Fatalf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", tunIP, subnetMask), "dev", tun.Name()); err != nil {
		log.Fatalf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(iranIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", iranIP), "dev", tun.Name()); err != nil {
			log.Fatalf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", iranIP), "dev", tun.Name()); err != nil {
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

func serverSide(tun *water.Interface, serverPort int, secretKey string, verbose bool) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		return fmt.Errorf("starting server failed: %v", err)
	}
	defer listener.Close()

	if verbose {
		log.Printf("Server listening on port %d\n", serverPort)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept any connection: %v", err)
			continue
		}
		go handle(conn, tun, secretKey, verbose)
	}
}

func handle(conn net.Conn, tun *water.Interface, secretKey string, verbose bool) {
	defer conn.Close()

	buf := make([]byte, len(secretKey))
	_, err := conn.Read(buf)
	if err != nil {
		if verbose {
			log.Printf("Reading authentication key failed: %v", err)
		}
		return
	}

	if string(buf) != secretKey {
		if verbose {
			log.Printf("Invalid secret key")
		}
		return
	}

	if verbose {
		log.Println("Client connected")
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromClient(conn, tun, clientToTun, verbose)
	go toTun(tun, clientToTun, verbose)
	go fromTun(tun, tunToClient, verbose)
	go toClient(conn, tunToClient, verbose)

	select {}
}

func fromClient(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			if verbose {
				log.Printf("Couldn't read packet length from client: %v", err)
			}
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := data(conn, buff)
		if err != nil {
			if verbose {
				log.Printf("Couldn't read data from client: %v", err)
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

func toClient(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		if _, err := conn.Write(packet); err != nil {
			if verbose {
				log.Printf("Couldn't write to client: %v", err)
			}
		}
	}
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

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func iPv6(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil
}
