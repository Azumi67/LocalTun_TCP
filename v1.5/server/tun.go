package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/Azumi67/LocalTun_TCP/v1.5/utils"
)

func SetupTun(config *Config, log *logrus.Logger) (*water.Interface, error) {
	cfg := water.Config{DeviceType: water.TUN}
	cfg.Name = config.TunName
	tun, err := water.New(cfg)
	if err != nil {
		return nil, err
	}

	if err := utils.Cmd(log, "ip", "link", "set", "dev", tun.Name(), "up"); err != nil {
		return nil, err
	}

	if err := utils.Cmd(log, "ip", "link", "set", "dev", tun.Name(), "mtu", fmt.Sprintf("%d", config.MTU)); err != nil {
		return nil, err
	}

	if err := utils.Cmd(log, "ip", "addr", "add", fmt.Sprintf("%s/%s", config.TunIP, config.SubnetMask), "dev", tun.Name()); err != nil {
		return nil, err
	}

	return tun, nil
}
