# cloudevent-gateway

A small service that ingests various `Event` sources, processes the Event based on the type, creates `CloudEvents` and sends the CloudEvents to a specified `Sink`. It is designed to be small and self-contained and depends on only corev1 `Kubernetes` resources to run.

The provided example manifests allow it to be started as a Deployment.

## Config

The following ENV variables exist to configure the gateway:

- EVENTNAMES - List of event names to be interested in, ";" delimited.
- EVENTTYPE - The type of source events the gateway should expect, defaults to `github`.
- EVENT_SINK - The http endpoint to send CloudEvents to.
- EVENT_SOURCE - The origin of the events. For github eventtype, this should be repo name in `iancoffey/cloudevent-gateway` format.

### Installation

You can apply the current latest example with `ko`.

`ko apply -f releases/`
