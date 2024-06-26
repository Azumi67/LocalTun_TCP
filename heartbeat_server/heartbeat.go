package heartbeatserver

import (
	"net"
	"time"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func trueHeartbeat(conn net.Conn, heartbeatInterval int) {
	heartbeatTiktok := time.NewTicker(time.Duration(heartbeatInterval) * time.Second)
	defer heartbeatTiktok.Stop()

	for range heartbeatTiktok.C {
		if _, err := conn.Write([]byte("ping")); err != nil {
			log.Warnf("Sending heartbeat failed: %v", err)
			conn.Close()
			return
		}
	}
}
