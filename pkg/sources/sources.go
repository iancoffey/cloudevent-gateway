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
	EventName string
	timeout   int
	hook
}

type Event string

type EventReceiver interface {
	New() *EventReceiver
	HandleEvent() error
}

func NewServer(eventType, events []Event, hook gh.Webhook) *EventServer {
	return &EventServer{
		EventType: eventType,
		Events:    eventName,
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
