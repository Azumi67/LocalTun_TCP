package workerserver

import (
	"crypto/sha256"
	"net"
	"sync"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/heartbeat_server"
	"github.com/Azumi67/LocalTun_TCP/nonce_server"
	"github.com/Azumi67/LocalTun_TCP/server"
)

var log = logrus.New()

func Worker(tun *water.Interface, workerChan chan net.Conn, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range workerChan {
		handleClient(conn, tun, secretKey, verbose, useSmux, enableHeartbeat, heartbeatInterval)
	}
}

func handleClient(conn net.Conn, tun *water.Interface, secretKey string, verbose bool, useSmux bool, enableHeartbeat bool, heartbeatInterval int) {
	defer conn.Close()

	nonce, err := nonceserver.UniqueNonce()
	if err != nil {
		log.Warnf("Generating nonce failed: %v", err)
		return
	}

	if _, err := conn.Write([]byte(nonce)); err != nil {
		log.Warnf("Sending unique nonce failed: %v", err)
		return
	}

	responseBuf := make([]byte, 32)
	if _, err := conn.Read(responseBuf); err != nil {
		log.Warnf("Reading the response wasn't possible: %v", err)
		return
	}

	hashkey := sha256.New()
	hashkey.Write([]byte(nonce))
	hashkey.Write([]byte(secretKey))
	expectThisHash := hashkey.Sum(nil)

	if !nonceserver.CalHashes(responseBuf, expectThisHash) {
		log.Warnf("Wrong authentication response")
		return
	}

	if enableHeartbeat {
		go heartbeatserver.trueHeartbeat(conn, heartbeatInterval)
	}

	if useSmux {
		serversmux.HandleSmux(conn, tun, verbose)
		return
	}

	clientToTun := make(chan []byte, 100)
	tunToClient := make(chan []byte, 100)

	go tunserver.FromClient(conn, tun, clientToTun, verbose)
	go tunserver.ToTun(tun, clientToTun, verbose)
	go tunserver.FromTun(tun, tunToClient, verbose)
	go tunserver.ToClient(conn, tunToClient, verbose)

	select {}
}
