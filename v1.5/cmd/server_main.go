package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/Azumi67/LocalTun_TCP/v1.5/server"
	"github.com/Azumi67/LocalTun_TCP/v1.5/utils"
	"net"
)

func main() {
	config := server.ParseConfig()

	log := utils.SetupLogger(config.Verbose, "/etc/server.log")

	tun, err := server.SetupTun(config, log)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
	if err != nil {
		log.Fatalf("Couldn't start the server: %v", err)
	}
	defer listener.Close()

	log.Infof("Server listening on port %d", config.ServerPort)

	for i := 0; i < config.NumWorkers; i++ {
		go server.Worker(log, listener, tun, config)
	}

	select {}
}
