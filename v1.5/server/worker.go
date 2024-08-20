package server

import (
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"net"
	"github.com/Azumi67/LocalTun_TCP/v1.5/utils"
)

func Worker(log *logrus.Logger, listener net.Listener, tun *water.Interface, config *Config) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Couldn't accept connection: %v", err)
			continue
		}

		if config.TCPNoDelay {
			utils.SetTCPNoDelay(conn, log)
		}

		log.Info("Client connected")

		go handle(log, conn, tun, config.SecretKey)
	}
}
