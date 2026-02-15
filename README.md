FREVOD is a flexible FFmpeg-based video encoding service designed for scalable Video On Demand (VOD) pipelines.

It abstracts FFmpeg complexity and integrates seamlessly with distributed systems, message queues, REST APIs, and cloud storage providers.

You can choose the format, either HLS or DASH.

Frevod has many ways to receive inputs and store outputs:

INPUTS:
- Your terminal / CLI app
- REST API
- RabbitMQ or Kafka queue
- Download from any source (as long as you provide the link)

OUTPUTS:
- Local storage
- Any storage compatible with S3 API

It can publish a message to a queue after finishing an encoding process.
It has built in monitoring and metrics with Prometheus.

This allows Frevod to fit in any architecture or pipeline.
Praise FFmpeg!
