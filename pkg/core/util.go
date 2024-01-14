package core

import (
	"net"

	"github.com/hashicorp/go-multierror"
)

// ParseIPsAndNetworks parses a list of strings into a list of net.IP and net.IPNet
func ParseIPsAndNetworks(list []string) ([]net.IP, []net.IPNet, error) {
	ips := make([]net.IP, 0, len(list))
	networks := make([]net.IPNet, 0, len(list))

	var errResult error

	for _, s := range list {
		ip := net.ParseIP(s)
		if ip != nil {
			ips = append(ips, ip)
			continue
		}

		_, ipNet, err := net.ParseCIDR(s)
		if err == nil {
			networks = append(networks, *ipNet)
			continue
		}

		errResult = multierror.Append(errResult, &net.ParseError{Type: "IP address or network", Text: s})
	}

	if len(ips) == 0 {
		ips = nil
	}
	if len(networks) == 0 {
		networks = nil
	}

	return ips, networks, errResult
}

// IsIpInListOfIpsOrListOfNetworks checks if an IP is in a list of IPs or a list of networks
func IsIpInListOfIpsOrListOfNetworks(ip net.IP, ips []net.IP, networks []net.IPNet) bool {
	for _, ip2 := range ips {
		if ip.Equal(ip2) {
			return true
		}
	}

	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
