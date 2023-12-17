package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/sirupsen/logrus"
	"github.com/things-go/go-socks5"
)

type Config struct {
	Port     int    `env:"PROXY_PORT" envDefault:"1080"`
	LogLevel string `env:"LogLevel" envDefault:"info"`
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
	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(logrus.StandardLogger()),
	)

	// Create SOCKS5 proxy on localhost port 8000
	logrus.WithField("port", config.Port).Info("Starting proxy...")
	if err := server.ListenAndServe("tcp", fmt.Sprintf(":%d", config.Port)); err != nil {
		return err
	}

	return nil
}
