package core

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/things-go/go-socks5"
)

func TestPolicy_Allow(t *testing.T) {
	_, classANet, err := net.ParseCIDR("1.0.0.0/8")
	require.NoError(t, err)
	_, classCNet, err := net.ParseCIDR("1.1.1.0/24")
	require.NoError(t, err)

	tests := []struct {
		name     string
		Rules    []Rule
		sourceIp string
		expected bool
	}{
		{
			name:     "empty rules",
			Rules:    []Rule{},
			sourceIp: "1.1.1.1",
			expected: false,
		},
		{
			name: "source ip matches single allow rule",
			Rules: []Rule{
				{
					From:   Source{Ips: []net.IP{net.ParseIP("1.1.1.1")}},
					Action: Allow,
				},
			},
			sourceIp: "1.1.1.1",
			expected: true,
		},
		{
			name: "source ip does not match single allow rule",
			Rules: []Rule{
				{
					From:   Source{Ips: []net.IP{net.ParseIP("1.1.1.1")}},
					Action: Allow,
				},
			},
			sourceIp: "3.3.3.3",
			expected: false,
		},
		{
			name: "allow rule matches first",
			Rules: []Rule{
				{
					From:   Source{CIDR: []net.IPNet{*classANet}},
					Action: Allow,
				},
				{
					From:   Source{CIDR: []net.IPNet{*classCNet}},
					Action: Deny,
				},
			},
			sourceIp: "1.1.1.1",
			expected: true,
		},
		{
			name: "deny rule matches first",
			Rules: []Rule{
				{
					From:   Source{CIDR: []net.IPNet{*classCNet}},
					Action: Deny,
				},
				{
					From:   Source{CIDR: []net.IPNet{*classANet}},
					Action: Allow,
				},
			},
			sourceIp: "1.1.1.1",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Policy{Rules: tt.Rules}
			req := socks5.Request{RemoteAddr: &net.IPAddr{IP: net.ParseIP(tt.sourceIp)}}

			_, allow := p.Allow(context.Background(), &req)
			assert.Equal(t, tt.expected, allow)
		})
	}
}
