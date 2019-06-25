# cloudevent-gateway

A small service that ingests various `Event` sources, processes the Event based on the type, creates `CloudEvents` and sends the CloudEvents to a specified `Sink`. It is designed to be small and self-contained and depends on only corev1 `Kubernetes` resources to run.

The provided example manifests allow it to be started as a Deployment.

## Config

The following ENV variables exist to configure the gateway:

EVENTNAME
EVENTTYPE
	EventType string   `env:"EVENTTYPE,default=github"`
	Events    []string `env:"EVENTNAME,default=push"`
	Secret    string   `env:"GITHUB_SECRET"`
	Sink      string   `env:"EVENT_SINK"`
	Source    string   `env:"EVENT_SOURCE"`


### Installation

You can apply the current latest example with `ko`.

`ko apply -f releases/`
