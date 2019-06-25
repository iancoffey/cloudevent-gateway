package github

import (
	"context"
	"fmt"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
	gh "gopkg.in/go-playground/webhooks.v5/github"
)

const (
	deliveryHeader = "X-GitHub-Delivery"
	eventHeader    = "X-GitHub-Event"
)

type GithubReceiver struct {
	sink     string
	logger   *log.Entry
	ghClient *gh.Webhook
	ceClient cloudevents.Client
	source   string // use for ownerRepo
}

func New(sink string, logger *log.Entry) (*GithubReceiver, error) {
	ghClient, err := gh.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create gh client: %q", err)
	}

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(sink),
		cloudevents.WithEncoding(cloudevents.HTTPBinaryV02), // TODO: make this a config var
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudEvent transport. error=%q", err)
	}
	ceClient, err := cloudevents.NewClient(t)
	if err != nil {
		return nil, fmt.Errorf("unable to create cloudevent client. error=%q", err)
	}

	return &GithubReceiver{
		sink:     sink,
		logger:   logger,
		ghClient: ghClient,
		ceClient: ceClient,
	}, nil
}

func (g *GithubReceiver) HandleEvent(events []string, r *http.Request) {

	eventID := g.ID(r)
	if eventID == "" {
		g.logger.Info("EventID not found - X-GitHub-Delivery unset")
		return
	}
	source := g.Source(r)
	if source == "" {
		g.logger.Info("Github Source not found - EVENT_SOURCE environment variable unset")
		return
	}
	eventName := g.EventName(r)
	if eventName == "" {
		g.logger.Info("Event Name not found - X-GitHub-Event unset")
	}
	if !eventInEvents(eventName, events) {
		g.logger.Fatalf("at=skipping-event event=%s", eventName)
		return
	}
	ghEvent, err := g.ghClient.Parse(r, stringsToEvents(events)...)
	if err != nil {
		if err == gh.ErrEventNotFound {
			g.logger.Info("Event not found")
			return
		}
	}

	event := cloudevents.NewEvent(cloudevents.VersionV02)
	event.SetID(eventID)
	event.SetType(eventName)
	event.SetSource(g.Source(r))
	event.SetExtension(eventHeader, eventName)
	event.SetExtension(deliveryHeader, eventID)
	event.SetSubject(eventName + "-" + eventID) // what should this be?
	event.SetData(ghEvent)

	if _, err := g.ceClient.Send(context.Background(), event); err != nil {
		g.logger.Infof("failed to send cloudevent error=%q", err)
		return
	}
}

func (g *GithubReceiver) ValidateType(r *http.Request) bool {
	event := g.EventName(r)
	if event == "" {
		return false
	}
	return true
}

func (g *GithubReceiver) EventName(r *http.Request) string {
	return r.Header.Get(eventHeader)
}

func (g *GithubReceiver) ID(r *http.Request) string {
	return r.Header.Get(deliveryHeader)
}

func (g *GithubReceiver) Source(r *http.Request) string {
	return fmt.Sprintf("github.com/%s", g.source)
}

func eventInEvents(event string, events []string) bool {
	for _, e := range events {
		if e == event {
			return true
		}
	}
	return false
}

func stringsToEvents(s []string) []gh.Event {
	c := make([]gh.Event, len(s))
	for i, v := range s {
		c[i] = gh.Event(v)
	}
	return c
}
