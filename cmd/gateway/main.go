package main

import (
	"github.com/iancoffey/cloudevent-gateway/pkg/sources"
	"github.com/iancoffey/cloudevent-gateway/pkg/sources/github"
	"github.com/iancoffey/cloudevent-gateway/pkg/types"

	"github.com/joeshaw/envdecode"
	log "github.com/sirupsen/logrus"
)

const (
	githubEventType = "github"
)

// Everything inline will be broken out :D

// Handle slack events in cool demo

// two these will need to be updated
// once we understand the structure this will be in
// - we will need to accept an array of eventtypes probably
type Config struct {
	EventType string   `env:"EVENTTYPE,default=github"`
	Events    []string `env:"EVENTNAME,default=push"`
	Sink      string   `env:"EVENT_SINK"`
	Source    string   `env:"EVENT_SOURCE"`
}

// we can process events like event.Process()

func main() {
	logger := log.WithFields(log.Fields{
		"app": "cloudevent-gateway",
	})
	logger.Info("Booting gateway!")

	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		logger.Fatal("Failed to decode env error=%q")
	}
	receivers := []types.EventReceiver{}

	switch cfg.EventType {
	case githubEventType:
		ghWebhook, err := github.New(cfg.Sink, logger)
		if err != nil {
			logger.Fatalf("Failed to create github client error=%q", err)
		}
		receivers = append(receivers, ghWebhook)
	default:
		logger.Fatal("Unrecognized Eventtype")
	}

	if len(receivers) == 0 {
		logger.Fatal("at=no-receivers-defined")
	}

	server := sources.NewServer(cfg.Events, receivers, logger)
	if err := server.Run(); err != nil {
		logger.Fatalf("Server Failed error=%q", err)
	}
}
