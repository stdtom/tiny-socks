package core

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIPsAndNetworks(t *testing.T) {
	tests := []struct {
		name           string
		list           []string
		expectIps      []net.IP
		expectNetworks []net.IPNet
		expectErr      assert.ErrorAssertionFunc
	}{
		{name: "empty list", list: []string{}, expectIps: nil, expectNetworks: nil, expectErr: assert.NoError},
		{name: "single ipv4", list: []string{"10.0.0.1"}, expectIps: []net.IP{net.IPv4(10, 0, 0, 1)}, expectNetworks: nil, expectErr: assert.NoError},
		{name: "single ipv6", list: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}, expectIps: []net.IP{net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")}, expectNetworks: nil, expectErr: assert.NoError},
		{name: "single ipv4 network", list: []string{"10.0.0.0/8"}, expectIps: nil, expectNetworks: []net.IPNet{{IP: net.IPv4(10, 0, 0, 0).To4(), Mask: net.CIDRMask(8, 32)}}, expectErr: assert.NoError},
		{name: "single ipv6 network", list: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64"}, expectIps: nil, expectNetworks: []net.IPNet{{IP: net.ParseIP("2001:0db8:85a3::"), Mask: net.CIDRMask(64, 128)}}, expectErr: assert.NoError},
		{name: "single invalid", list: []string{"1234"}, expectIps: nil, expectNetworks: nil, expectErr: assert.Error},
		{name: "single invalid network", list: []string{"1234/8"}, expectIps: nil, expectNetworks: nil, expectErr: assert.Error},
		{name: "single invalid ip", list: []string{"1234.5678.9012.3456"}, expectIps: nil, expectNetworks: nil, expectErr: assert.Error},
		{name: "single invalid ip net", list: []string{"1234.5678.9012.3456/8"}, expectIps: nil, expectNetworks: nil, expectErr: assert.Error},
		{
			name:           "mixed",
			list:           []string{"1.1.1.1", "10.0.0.0/8", "2.2.2.2", "192.168.200.0/24"},
			expectIps:      []net.IP{net.IPv4(1, 1, 1, 1), net.IPv4(2, 2, 2, 2)},
			expectNetworks: []net.IPNet{{IP: net.IPv4(10, 0, 0, 0).To4(), Mask: net.CIDRMask(8, 32)}, {IP: net.IPv4(192, 168, 200, 0).To4(), Mask: net.CIDRMask(24, 32)}},
			expectErr:      assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ips, networks, err := ParseIPsAndNetworks(tt.list)
			tt.expectErr(t, err)
			assert.Equal(t, tt.expectIps, ips)
			assert.Equal(t, tt.expectNetworks, networks)

		})
	}
}

func TestIsIpInListOfIpsOrListOfNetworks(t *testing.T) {
	tests := []struct {
		name     string
		ip       net.IP
		ips      []string
		networks []string
		expected bool
	}{
		{"empty lists", net.IPv4(10, 0, 0, 1), []string{}, []string{}, false},
		{"empty ip", net.IP{}, []string{}, []string{}, false},
		{"ip in list of ips", net.IPv4(1, 1, 1, 1), []string{"1.1.1.1", "8.8.8.8"}, []string{"10.0.0.0/8", "192.168.0.0/16"}, true},
		{"ip in list of ips", net.IPv4(8, 8, 8, 8), []string{"1.1.1.1", "8.8.8.8"}, []string{"10.0.0.0/8", "192.168.0.0/16"}, true},
		{"ip in list of networks", net.IPv4(10, 10, 10, 10), []string{"1.1.1.1", "8.8.8.8"}, []string{"10.0.0.0/8", "192.168.0.0/16"}, true},
		{"ip in list of networks", net.IPv4(192, 168, 254, 254), []string{"1.1.1.1", "8.8.8.8"}, []string{"10.0.0.0/8", "192.168.0.0/16"}, true},
		{"ip not in lists", net.IPv4(212, 212, 212, 212), []string{"1.1.1.1", "8.8.8.8"}, []string{"10.0.0.0/8", "192.168.0.0/16"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// join ips and networks
			list := append(tt.ips, tt.networks...)

			ips, networks, err := ParseIPsAndNetworks(list)
			assert.NoError(t, err)

			got := IsIpInListOfIpsOrListOfNetworks(tt.ip, ips, networks)
			assert.Equal(t, tt.expected, got)
		})
	}
}
