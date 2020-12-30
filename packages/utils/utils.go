package utils

import (
	"errors"
	"net"
)

// GetIPv4Address ... returns IPv4 address
func GetIPv4Address() (net.IP, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return []byte{}, err
	}
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipV4 := ipNet.IP.To4(); ipV4 != nil {
				return ipV4, nil
			}
		}
	}
	return []byte{}, errors.New("Cannot find IPv4 address")
}
