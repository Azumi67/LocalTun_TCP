package tunclient

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"sync"
        "github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

var log = logrus.New()

func tunUp(tunName string, clientTunIP string, serverTunIP string, subnetMask string, mtu int) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = tunName
	tun, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		return nil, fmt.Errorf("Couldn't bring up the TUN device: %v", err)
	}

	if err := cmd("ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		return nil, fmt.Errorf("Couldn't set MTU: %v", err)
	}

	if err := cmd("ip", "addr", "add", fmt.Sprintf("%s/%s", clientTunIP, subnetMask), "dev", tun.Name()); err != nil {
		return nil, fmt.Errorf("Couldn't assign private IP address to TUN device: %v", err)
	}

	if iPv6(serverTunIP) {
		if err := cmd("ip", "-6", "route", "add", fmt.Sprintf("%s/128", serverTunIP), "dev", tun.Name()); err != nil {
			return nil, fmt.Errorf("Adding route for private IPv6 failed: %v", err)
		}
	} else {
		if err := cmd("ip", "route", "add", fmt.Sprintf("%s/32", serverTunIP), "dev", tun.Name()); err != nil {
			return nil, fmt.Errorf("Adding route for private IPv4 failed: %v", err)
		}
	}

	return tun, nil
}

func HandleServer(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, stop chan struct{}) {
	if _, err := conn.Write([]byte(secretKey)); err != nil {
		log.Warnf("Couldn't send authentication key: %v", err)
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fromServer(conn, tun, clientToTun, verbose, stop)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		toTun(tun, clientToTun, verbose, stop)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fromTun(tun, tunToClient, verbose, stop)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		toServer(conn, tunToClient, verbose, stop)
	}()

	wg.Wait()
}

func fromServer(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			pcktLength := make([]byte, 2)
			if _, err := conn.Read(pcktLength); err != nil {
				log.Warnf("Couldn't read packet length from server: %v", err)
				close(clientToTun)
				return
			}

			length := binary.BigEndian.Uint16(pcktLength)
			buff := make([]byte, length)

			_, err := data(conn, buff)
			if err != nil {
				log.Warnf("Couldn't read data from server: %v", err)
				close(clientToTun)
				return
			}

			clientToTun <- buff
		}
	}
}

func toTun(tun *water.Interface, clientToTun chan []byte, verbose bool, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case buff := <-clientToTun:
			if _, err := tun.Write(buff); err != nil {
				log.Warnf("Couldn't write to TUN device: %v", err)
			}
		}
	}
}

func fromTun(tun *water.Interface, tunToClient chan []byte, verbose bool, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			buff := make([]byte, 1500)
			n, err := tun.Read(buff)
			if err != nil {
				log.Warnf("Couldn't read from TUN device: %v", err)
				continue
			}

			packet := make([]byte, 2+n)
			binary.BigEndian.PutUint16(packet[:2], uint16(n))
			copy(packet[2:], buff[:n])

			tunToClient <- packet
		}
	}
}

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case packet := <-tunToClient:
			if _, err := conn.Write(packet); err != nil {
				log.Warnf("Couldn't write to server: %v", err)
				return
			}
		}
	}
}

func data(conn net.Conn, buff []byte) (int, error) {
	n, err := conn.Read(buff)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func iPv6(address string) bool {
	return net.ParseIP(address).To4() == nil
}
