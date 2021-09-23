# OpenTelemetry Workshop for StrangeLoop 2021

Hello and welcome to the OpenTelemetry Workshop! Here you will learn how to:

- Automatically instrumenting an application written in Go, Java, Node, or Python
  with spans and sending the data to an APM tool like Jaeger.
- Processing data including CRUD metadata operations via the OpenTelemetry
  Collector.
- Collecting host metrics and sending the data to a metrics tool like
  Prometheus.
- Ensure your telemetry data adheres to open-standards and remains
  vendor-agnostic.

## Prerequisites

- Docker with `docker-compose`

Optionally, it would be good to review [OpenTelemety
concepts](https://opentelemetry.io/docs/concepts/).

## Getting Started

```
$ cd docker
$ docker-compose up --build
```

- Go is listening outside docker on `7070`.
- Java is listening outside docker on `8080`.
- Node is listening outside docker on `3000`.
- Python is listening outside docker on `5000`.

Exercise the API:

```
$ curl -i localhost:3000/hello
$ curl -i localhost:8080/hello/proxy/python/node/java
```

Once instrumented for OpenTelemetry, view traces in Jaeger:
http://localhost:16686. Also note that an OpenTelemetry Collector is available
and can receive data via:

- Jaeger inside docker on `14250` and `14268`
- OTLP inside docker on `4317` and `4318`
- Zipkin inside/outside docker on `9411`

A Prometheus server which receives OpenTelemetry Collector metrics is available at http://localhost:9090.

## Labs

- Lab 101: Complete one or more of the following
  - Instrument a Go application
  - Instrument a Java application
  - Instrument a NodeJS application
  - Instrument a Python application
- Lab 102: Reconfigure applications to leverage the OpenTelemetry Collector
- Lab 201: Processing trace data with the OpenTelemetry Collector
- Lab 202: Collecting and processing metric data with the OpenTelemetry Collector
- Lab 203: Collecting and processing log data with the OpenTelemetry Collector
