package monitorclient

import (
	"os/exec"
	"time"
	"net"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func monitor(serverTunIP string, pingInterval int, serviceName string) {
	for {
		time.Sleep(time.Duration(pingInterval) * time.Second)
		err := ping(serverTunIP)
		if err != nil {
			log.Warnf("Ping to Server's PrivateIP failed: %v", err)
			restart(serviceName)
		}
	}
}

func ping(serverTunIP string) error {
	var cmd *exec.Cmd
	if iPv6(serverTunIP) {
		cmd = exec.Command("ping6", "-c", "1", serverTunIP)
	} else {
		cmd = exec.Command("ping", "-c", "1", serverTunIP)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Warnf("Ping command output: %s", out)
	}
	return err
}

func restart(serviceName string) {
	log.Warn("Restarting client..")
	cmd := exec.Command("systemctl", "restart", serviceName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Restarting client failed: %v", err)
	}
}

func iPv6(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil
}
