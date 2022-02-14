package network

import (
	"net"
)

// GetLocalIP returns your local IP address to use while
// pretending to be bank.com.
// This function may change in grading; do not hardcode this functionality.
func GetLocalIP() net.IP {
	inter, _ := net.InterfaceByName("eth0")
	adList, _ := inter.Addrs()
	ip, _, _ := net.ParseCIDR(adList[0].String())
	return ip
}

// GetBankIP returns the true IP address of bank.com.
// This function may change in grading; do not hardcode this address.
func GetBankIP() net.IP {
	return net.ParseIP("10.38.8.3")
}
