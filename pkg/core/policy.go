package core

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/things-go/go-socks5"
)

type Policy struct {
	Rules []Rule
}

// Allow implements RuleSet.Allow by checking if the request's source IP is in the list of allowed IPs
func (p *Policy) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	logrus.WithField("remote_addr", req.RemoteAddr).WithField("local_addr", req.LocalAddr).WithField("dest_addr", req.DestAddr).Debug("policy.allow called")
	for _, rule := range p.Rules {
		_, action := rule.evaluate(ctx, req)
		switch action {
		case Deny:
			return ctx, false
		case Allow:
			return ctx, true
		}
	}
	return ctx, false
}
