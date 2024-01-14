package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/sirupsen/logrus"
	"github.com/stdtom/tiny-proxy/pkg/core"
	"github.com/things-go/go-socks5"
)

type Config struct {
	IpAddress     string   `env:"PROXY_IP" envDefault:"127.0.0.1"`
	Port          int      `env:"PROXY_PORT" envDefault:"1080"`
	AllowedSource []string `env:"ALLOWED_SOURCE" envDefault:"127.0.0.1"`
	LogLevel      string   `env:"LogLevel" envDefault:"info"`
}

func main() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stderr)

	var config Config
	if err := env.Parse(&config); err != nil {
		logrus.WithError(err).Panic("could not parse env config")
	}

	l, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)

	if err := run(config); err != nil {
		logrus.WithError(err).Error("Exiting on error")
		os.Exit(1)
	}
}

func run(config Config) error {
	ipList, networkList, err := core.ParseIPsAndNetworks(config.AllowedSource)
	if err != nil {
		return err
	}

	ruleAllowedSources := core.Rule{
		From: core.Source{
			Ips:  ipList,
			CIDR: networkList,
		},
		Action: core.Allow,
	}

	policy := core.Policy{Rules: []core.Rule{ruleAllowedSources}}
	logrus.WithField("policy", policy).Info("Policy created")

	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(logrus.StandardLogger()),
		socks5.WithRule(&policy),
	)

	// Start listening on defined ip address and port
	logrus.WithField("port", config.Port).Info("Starting proxy...")
	if err := server.ListenAndServe("tcp", fmt.Sprintf("%s:%d", config.IpAddress, config.Port)); err != nil {
		return err
	}

	return nil
}
