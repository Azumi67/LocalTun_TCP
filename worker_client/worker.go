package workerclient

import (
	"net"
	"sync"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/smux_client"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_client"
	"github.com/Azumi67/LocalTun_TCP/hash"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Worker(tun *water.Interface, workerOnechan chan net.Conn, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range workerOnechan {
		handleServer(conn, tun, secretKey, verbose, useSmux, enableHeartbeat, heartbeatInterval)
	}
}

func handleServer(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int) {
	challenge := make([]byte, 32)
	if _, err := conn.Read(challenge); err != nil {
		log.Warnf("Couldn't read challenge from server: %v", err)
		return
	}

	hashedKey := hash.GenHash(string(challenge), secretKey)

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

	go smuxclient.FromServer(conn, tun, clientToTun, verbose)
	go smuxclient.ToTun(tun, clientToTun, verbose)
	go smuxclient.FromTun(tun, tunToClient, verbose)
	go smuxclient.ToServer(conn, tunToClient, verbose)

	select {}
}
