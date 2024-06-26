package monitorserver

import (
	"os/exec"
	"time"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Monitor(clientTunIP string, pingInterval int, serviceName string) {
	for {
		time.Sleep(time.Duration(pingInterval) * time.Second)
		if err := cmd("ping", "-c", "1", "-W", "2", clientTunIP); err != nil {
			log.Infof("ping to %s failed, restarting the service: %s", clientTunIP, serviceName)
			if err := cmd("systemctl", "restart", serviceName); err != nil {
				log.Errorf("failed to restart the service: %s", serviceName)
			}
		} else {
			//log.Infof("ping to %s succeeded", clientTunIP)
		}
	}
}

func cmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
