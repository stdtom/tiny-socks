package core

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/things-go/go-socks5"
)

type Source struct {
	Ips     []net.IP
	NotIps  []net.IP
	CIDR    []net.IPNet
	NotCIDR []net.IPNet
}

type Rule struct {
	From   Source
	Action Action
}

func (r Rule) evaluate(ctx context.Context, req *socks5.Request) (context.Context, Action) {
	var srcIpAddr net.IPAddr

	switch srcAddr := req.RemoteAddr.(type) {
	case *net.IPAddr:
		srcIpAddr = *srcAddr
	case *net.TCPAddr:
		srcIpAddr = net.IPAddr{IP: srcAddr.IP}
	}

	ruleMatches := r.matchesSource(srcIpAddr.IP)
	if ruleMatches {
		logrus.WithField("rule", r).WithField("source_ip", srcIpAddr.IP).WithField("action", r.Action.String()).Debug("rule matched")
		return ctx, r.Action
	}
	logrus.WithField("rule", r).WithField("source_ip", srcIpAddr.IP).WithField("action", Unknown.String()).Debug("rule did not match")

	return ctx, Unknown
}

func (r Rule) matchesSource(ip net.IP) bool {
	isPositiveListDefined := (len(r.From.Ips) + len(r.From.CIDR)) > 0
	isNegativeListDefined := (len(r.From.NotIps) + len(r.From.NotCIDR)) > 0

	isInPositiveList := IsIpInListOfIpsOrListOfNetworks(ip, r.From.Ips, r.From.CIDR)
	isInNegativeList := IsIpInListOfIpsOrListOfNetworks(ip, r.From.NotIps, r.From.NotCIDR)

	switch {
	case isPositiveListDefined && !isNegativeListDefined:
		return isInPositiveList

	case !isPositiveListDefined && isNegativeListDefined:
		return !isInNegativeList

	case isPositiveListDefined && isNegativeListDefined:
		return isInPositiveList && !isInNegativeList

	case !isPositiveListDefined && !isNegativeListDefined:
		return false
	}

	return false
}
