package main

import (
	"github.com/iancoffey/cloudevent-gateway/pkg/sources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	envSecret       = "GITHUB_SECRET"
	githubEventType = "github"
)

// Everything inline will be broken out :D

// Handle slack events in cool demo

// two these will need to be updated
// once we understand the structure this will be in
// - we will need to accept an array of eventtypes probably
type Config struct {
	EventType string          `env:"EVENTTYPE,default=github"`
	Events    []sources.Event `env:"EVENTNAME,default=push"`
	Secret    string          `env:"GITHUB_SECRET"`
}

// we can process events like event.Process()

func main() {
	logger := log.WithFields(log.Fields{
		"app": "cloudevent-gateway",
	})
	logger.Info("Booting gateway!")

	var cfg Config
	err := envdecode.Decode(&cfg)

	switch EventType {
	case githubEventType:
		webhook, err := gh.New(gh.Options.Secret(cfg.Secret))
		if err != nil {
			logger.Fatalf("at=github-new", err)
		}
	default:
		logger.Fatal("at=unrecognized-eventtype")
	}

	server := sources.NewServer(cfg.EventType, cfg.EventName, webhook)
	if err := server.Run(); err != nil {
		logger.Fatalf("at=server-failed", err)
	}
}
