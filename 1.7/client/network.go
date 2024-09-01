package utils

import "net"

func IPv6(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() == nil
}

func DefaultSubnet(ip string) string {
	if IPv6(ip) {
		return "64"
	}
	return "24"
}

