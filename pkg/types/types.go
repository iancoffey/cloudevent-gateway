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
	HandleEvent(string, []string, *http.Request)
	// Each event types provides its own mechanism for determining if it is a valid event.
	// - if the event received matches no known type, then we drop, log and continue
	// - if the event matches a type, but not what wqe are configured for, log and continue
	// - if we get the right type and it matches, process and send to Sink.
	ValidateType(*http.Request) (string, bool)
	// Return an ID for this type of event
	ID(r *http.Request) string
	// The source of the event, eg, for github it is the owner repository "iancoffey/cloudevent-gateway"
	Source(r *http.Request) string
}

func (e *EventServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, receiver := range e.Receivers {
			if eventName, ok := receiver.ValidateType(r); ok {
				receiver.HandleEvent(eventName, e.Events, r)
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
