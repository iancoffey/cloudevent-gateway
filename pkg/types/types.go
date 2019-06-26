package types

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	port = 8089
)

// We will ingest a github event and produce a cloudevent
type EventServer struct {
	Events    []string
	Receivers []EventReceiver
	Timeout   int // TODO implement timeouts
	Logger    *log.Entry
}

type EventReceiver interface {
	// Handle this specific eventtype, eg github, gitlab, slack, w/e
	HandleEvent([]string, *http.Request)
	// Each event types provides its own mechanism for determining if it is a valid event.
	ValidateType(*http.Request) bool
	// Return an ID for this type of event
	ID(r *http.Request) string
	// The source of the event, eg, for github it is the owner repository "iancoffey/cloudevent-gateway"
	Source(r *http.Request) string
}

func (e *EventServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, receiver := range e.Receivers {
			if ok := receiver.ValidateType(r); ok {
				receiver.HandleEvent(e.Events, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)

		e.Logger.Info("Receiver for event-type not found")
		return
	})

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}
