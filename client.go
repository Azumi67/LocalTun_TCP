package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"sync"
	"time"

	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_client"
	"github.com/Azumi67/LocalTun_TCP/smux_client"
	"github.com/Azumi67/LocalTun_TCP/client"
)

var log = logrus.New()

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g : 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")

	flag.Parse()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = *tunName
	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	tunclient.TunUp(tun, *tunIP, *serverTunIP, *subnetMask, *mtu, *tunName)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *kharejAddr, *serverPort))
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err)
	}
	defer conn.Close()

	log.Info("Connected to server")

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Warn("Converting to TCPConn failed")
		return
	}

	if err := tcpnodelayclient.noDelay(tcpConn); err != nil {
		log.Warnf("Setting up TCP no delay failed: %v", err)
		return
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if *useSmux {
			smuxclient.HandleSmux(conn, tun, *verbose, stop)
		} else {
			handleServer(conn, tun, *secretKey, *verbose, stop)
		}
	}()

	for {
		time.Sleep(10 * time.Second) 
		if !ping(*serverTunIP) {
			log.Warn("Server is unreachable, azumi is reconnecting..")
			close(stop) 
			conn.Close()
			wg.Wait() 
			break
		}
	}
}

func handleServer(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, stop chan struct{}) {
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
	return conn.Read(buff)
}

func ping(ip string) bool {
	// later it will be added or in script i will make a keepalive ping service
	return true
}
