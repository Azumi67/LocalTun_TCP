package main

import (
	"fmt"
	"net"
	"github.com/sirupsen/logrus"
	"github.com/Azumi67/LocalTun_TCP/v1.5/client"
	"github.com/Azumi67/LocalTun_TCP/v1.5/utils"
)

func main() {
	config := client.ParseConfig()

	log := utils.SetupLogger(config.Verbose, "/etc/client.log")

	tun, err := client.SetupTun(config, log)
	if err != nil {
		log.Fatalf("Couldn't create TUN device: %v", err)
	}
	defer tun.Close()

	conn, err := net.Dial("tcp", fmt.Sprintf("[%s]:%d", config.ServerAddr, config.ServerPort))
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err)
	}
	defer conn.Close()

	log.Info("Connected to server")

	if config.TCPNoDelay {
		utils.SetTCPNoDelay(conn, log)
	}

	for i := 0; i < config.NumWorkers; i++ {
		go client.Worker(log, conn, tun, config)
	}

	select {}
}
