package main

import (
	"flag"
	"log"
	"time"

	"github.com/Azumi67/LocalTun_TCP/v1.6/server/internal/config"
	"github.com/Azumi67/LocalTun_TCP/v1.6/server/internal/crypto"
	"github.com/Azumi67/LocalTun_TCP/v1.6/server/internal/network"
	"github.com/songgao/water"
	"github.com/xtaci/smux"
)

func main() {
	cfg := config.LoadServerConfig()

	publicKey, err := crypto.LoadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("Loading public key was not successful: %v", err)
	}

	tun, err := network.SetupTunDevice(cfg)
	if err != nil {
		log.Fatalf("Couldn't set up TUN device: %v", err)
	}
	defer tun.Close()

	listener, err := network.SetupListener(cfg.ServerPort)
	if err != nil {
		log.Fatalf("Starting server failed: %v", err)
	}
	defer listener.Close()

	log.Printf("Server listening on port %d\n", cfg.ServerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Couldn't accept connection: %v", err)
			continue
		}

		network.ConfigureConnection(conn, cfg)

		if cfg.EnableSmux {
			session, err := smux.Server(conn, nil)
			if err != nil {
				log.Printf("Couldn't create smux session: %v", err)
				conn.Close()
				continue
			}
			go network.HandleSmux(session, tun, publicKey, cfg.Verbose)
		} else {
			go network.Handle(conn, tun, publicKey, cfg.Verbose)
		}
	}
}
