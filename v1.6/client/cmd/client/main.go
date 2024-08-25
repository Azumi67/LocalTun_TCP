package main

import (
	"flag"
	"log"

	"github.com/Azumi67/LocalTun_TCP/v1.6/client/internal/config"
	"github.com/Azumi67/LocalTun_TCP/v1.6/client/internal/crypto"
	"github.com/Azumi67/LocalTun_TCP/v1.6/client/internal/network"
	"github.com/xtaci/smux"
)

func main() {
	cfg := config.LoadClientConfig()

	privateKey, err := crypto.LoadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Loading private key was not successful: %v", err)
	}

	tun, err := network.SetupTunDevice(cfg)
	if err != nil {
		log.Fatalf("Couldn't set up TUN device: %v", err)
	}
	defer tun.Close()

	conn, err := network.SetupConnection(cfg)
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err)
	}
	defer conn.Close()

	network.ConfigureConnection(conn, cfg)

	if cfg.EnableSmux {
		smuxSession, err := smux.Client(conn, nil)
		if err != nil {
			log.Fatalf("Couldn't create smux session: %v", err)
		}
		defer smuxSession.Close()

		conn, err = smuxSession.OpenStream()
		if err != nil {
			log.Fatalf("Opening smux stream failed: %v", err)
		}
		defer conn.Close()
	}

	if err := network.AuthenticateClient(conn, privateKey); err != nil {
		log.Fatalf("Client authentication failed: %v", err)
	}

	log.Println("Client authenticated")

	go network.HandleConnection(conn, tun, cfg.Verbose)

	network.ForwardTraffic(tun, conn, cfg.Verbose)
}
