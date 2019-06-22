package sources

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	port = 8089
)

// What sort of event are we listening for?
// We will ingest a github event and produce a cloudevent
type EventServer struct {
	EventType string
	Events    []Event
	Sink      string
	Receiver  *EventReceiver
	timeout   int
	hook
}

type Event string

type EventReceiver interface {
	// handle this specific eventtype, eg github, gitlab, slack, w/e
	HandleEvent() error
	// every Receiver type will be cycled through to determine the type
	// we will check for github headers to make that cal for github
	// if the event received matches no known type, then we drop, log and continue
	// if the event matches a type, but not what wqe are configured for, log and continue
	// if we get the right type and it matches, process and send to Sink.
	ValidateType(body, headers string) bool
}

func NewServer(eventType, events []Event, receiver *EventReceiver) *EventServer {
	return &EventServer{
		EventType: eventType,
		Events:    events,
		Receiver:  receiver,
	}
}

func (e *EventServer) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		event, err := hook.Parse(r, validEvents...)
		if err != nil {
			if err == gh.ErrEventNotFound {
				w.WriteHeader(http.StatusNotFound)

				log.Print("Event not found")
				return
			}
		}

		ra.HandleEvent(event, r.Header)
	})
	addr := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(addr, nil)
}
