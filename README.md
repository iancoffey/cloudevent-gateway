# cloudevent-gateway

A small service that ingests various `Event` sources, processes the Event based on the type, creates `CloudEvents` and sends the CloudEvents to a specified `Sink`. It is designed to be small and self-contained and depends on only corev1 `Kubernetes` resources to run.

The provided example manifests allow it to be started as a Deployment.
