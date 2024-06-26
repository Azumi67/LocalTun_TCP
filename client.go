package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
	"github.com/Azumi67/LocalTun_TCP/client"
	"github.com/Azumi67/LocalTun_TCP/smux_client"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_client"
	"github.com/Azumi67/LocalTun_TCP/monitor_client"
	"github.com/Azumi67/LocalTun_TCP/hash"
	"github.com/Azumi67/LocalTun_TCP/tcp_no_delay_client"
)

var log = logrus.New()

func main() {
	kharejAddr := flag.String("server-addr", "SERVER_IP", "Server Public IP address")
	serverPort := flag.Int("server-port", 800, "Server port")
	tunIP := flag.String("client-private", "2001:db8::2", "Client Private IP address")
	serverTunIP := flag.String("server-private", "2001:db8::1", "Server Private IP address")
	subnetMask := flag.String("subnet", "64", "Subnet mask (e.g: 24 or 64)")
	tunName := flag.String("device", "tun2", "TUN device name")
	secretKey := flag.String("key", "azumi", "Secret key for authentication")
	mtu := flag.Int("mtu", 1480, "MTU for TUN device")
	verbose := flag.Bool("verbose", false, "Enable logging")
	useSmux := flag.Bool("smux", false, "Enable smux multiplexing")
	enableHeartbeat := flag.Bool("heartbeat", false, "Enable heartbeat")
	heartbeatInterval := flag.Int("heartbeat-interval", 30, "Heartbeat interval in seconds")
	pingInterval := flag.Int("ping-interval", 10, "Ping interval in seconds")
	tcpNoDelay := flag.Bool("tcp-nodelay", false, "Enable TCP_NODELAY")
	serviceName := flag.String("service-name", "azumilocal", "name of the service to restart upon failure")

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

	go monitorclient.monitor(*serverTunIP, *pingInterval, *serviceName)

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *kharejAddr, *serverPort))
		if err != nil {
			log.Warnf("Couldn't connect to server: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if *tcpNoDelay {
			tcpnodelayclient.SetTCPNoDelay(conn)
		}

		defer conn.Close()

		log.Info("Connected to server")
		handleServer(conn, tun, *secretKey, *verbose, *useSmux, *enableHeartbeat, *heartbeatInterval)
	}
}

func handleServer(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int) {
	challenge := make([]byte, 32)
	if _, err := conn.Read(challenge); err != nil {
		log.Warnf("Couldn't read challenge from server: %v", err)
		return
	}

	hashedKey := hashclient.GenHash(challenge, secretKey)

	if _, err := conn.Write(hashedKey); err != nil {
		log.Warnf("Couldn't send auth response: %v", err)
		return
	}

	if enableHeartbeat {
		go heartbeatclient.trueHeartbeat(conn, heartbeatInterval)
	}

	if useSmux {
		smuxclient.HandleSmux(conn, tun, verbose)
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go fromServer(conn, tun, clientToTun, verbose)
	go tunclient.ToTun(tun, clientToTun, verbose)
	go tunclient.FromTun(tun, tunToClient, verbose)
	go toServer(conn, tunToClient, verbose)

	select {}
}

func fromServer(conn net.Conn, tun *water.Interface, clientToTun chan []byte, verbose bool) {
	for {
		pcktLength := make([]byte, 2)
		if _, err := conn.Read(pcktLength); err != nil {
			log.Warnf("Couldn't read packet length from server: %v", err)
			close(clientToTun)
			return
		}

		length := binary.BigEndian.Uint16(pcktLength)
		buff := make([]byte, length)

		_, err := tunclient.Data(conn, buff)
		if err != nil {
			log.Warnf("Couldn't read data from server: %v", err)
			close(clientToTun)
			return
		}

		clientToTun <- buff
	}
}

func toServer(conn net.Conn, tunToClient chan []byte, verbose bool) {
	for packet := range tunToClient {
		_, err := conn.Write(packet)
		if err != nil {
			log.Warnf("Couldn't write to server: %v", err)
		}
	}
}
