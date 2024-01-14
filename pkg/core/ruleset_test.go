package core

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/things-go/go-socks5"
)

func TestRule_evaluate(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		sourceIp string
		expected Action
	}{
		{
			name:     "unknown by default",
			rule:     Rule{},
			sourceIp: "1.1.1.1",
			expected: Unknown,
		},
		{
			name: "match source ips",
			rule: Rule{
				From: Source{
					Ips: []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")},
				},
				Action: Allow,
			},
			sourceIp: "2.2.2.2",
			expected: Allow,
		},
		{
			name: "not match source ips",
			rule: Rule{
				From: Source{
					Ips: []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")},
				},
				Action: Allow,
			},
			sourceIp: "3.3.3.3",
			expected: Unknown,
		},
		{
			name: "match source notIps",
			rule: Rule{
				From: Source{
					NotIps: []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")},
				},
				Action: Allow,
			},
			sourceIp: "3.3.3.3",
			expected: Allow,
		},
		{
			name: "not match source notIps",
			rule: Rule{
				From: Source{
					NotIps: []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2")},
				},
				Action: Allow,
			},
			sourceIp: "1.1.1.1",
			expected: Unknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := socks5.Request{RemoteAddr: &net.IPAddr{IP: net.ParseIP(tt.sourceIp)}}

			_, got := tt.rule.evaluate(context.Background(), &req)

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestRule_matchesSource(t *testing.T) {
	_, posNets, _ := ParseIPsAndNetworks([]string{"10.0.0.0/8"})
	posIps, _, _ := ParseIPsAndNetworks([]string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4"})
	_, negNets, _ := ParseIPsAndNetworks([]string{"10.1.0.0/16"})
	negIps, _, _ := ParseIPsAndNetworks([]string{"1.1.1.1", "2.2.2.2"})

	posSource := Source{
		Ips:  posIps,
		CIDR: posNets,
	}
	negSource := Source{
		NotIps:  negIps,
		NotCIDR: negNets,
	}
	combinedSource := Source{
		Ips:     posIps,
		CIDR:    posNets,
		NotIps:  negIps,
		NotCIDR: negNets,
	}
	noSource := Source{}

	tests := []struct {
		name     string
		source   Source
		ip       string
		expected bool
	}{
		{name: "in positive list, no negative list", source: posSource, ip: "10.10.10.10", expected: true},
		{name: "in positive list, no negative list", source: posSource, ip: "1.1.1.1", expected: true},
		{name: "not in positive list, no negative list", source: posSource, ip: "8.8.8.8", expected: false},
		{name: "in negative list, no positive list", source: negSource, ip: "10.1.2.3", expected: false},
		{name: "in negative list, no positive list", source: negSource, ip: "1.1.1.1", expected: false},
		{name: "not in negative list, no positive list", source: negSource, ip: "8.8.8.8", expected: true},
		{name: "not in negative list, no positive list", source: negSource, ip: "10.10.10.10", expected: true},
		{name: "in positive list, not in negative list", source: combinedSource, ip: "10.10.10.10", expected: true},
		{name: "in positive list, not in negative list", source: combinedSource, ip: "4.4.4.4", expected: true},
		{name: "in positive list and in negative list", source: combinedSource, ip: "1.1.1.1", expected: false},
		{name: "in positive list and in negative list", source: combinedSource, ip: "10.1.2.3", expected: false},
		{name: "neither in positive nor in negative list", source: combinedSource, ip: "8.8.8.8", expected: false},
		{name: "no positive list and no negative list", source: noSource, ip: "1.1.1.1", expected: false},
		{name: "no positive list and no negative list", source: noSource, ip: "3.3.3.3", expected: false},
		{name: "no positive list and no negative list", source: noSource, ip: "10.1.2.3", expected: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rule{
				From: tt.source,
			}
			ip := net.ParseIP(tt.ip)

			got := r.matchesSource(ip)

			assert.Equal(t, tt.expected, got)
		})
	}
}
