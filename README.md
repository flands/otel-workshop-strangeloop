# OpenTelemetry Workshop for StrangeLoop 2021

Hello and welcome to the OpenTelemetry Workshop! Here you will learn how to:

- Automatically instrumenting an application written in Java, Python, Node, or Go.
- Send span data to an APM tool like Jaeger.
- Processing data including CRUD metadata operations via the OpenTelemetry
  Collector.
- Collecting host metrics and sending the data to a metrics tool like
  Prometheus.
- Ensure your telemetry data adheres to open-standards and remains
  vendor-agnostic.

## Prerequisites

- Docker with `docker-compose`

Optionally, it would be good to review [OpenTelemetry
concepts](https://opentelemetry.io/docs/concepts/).

## Getting Started

```bash
cd docker
docker-compose up --build
```

- Java is listening outside docker on `8080`.
- Python is listening outside docker on `5000`.
- Node is listening outside docker on `3000`.
- Go is listening outside docker on `7070`.

Exercise the API:

```bash
curl -i localhost:3000/hello
curl -i localhost:8080/hello/proxy/python/node/java
```

Other included components:

- Once instrumented with OpenTelemetry, view traces in Jaeger:
http://localhost:16686.
- An OpenTelemetry Collector is available and can receive data via:
  - Jaeger inside docker on `14250` and `14268`
  - OTLP inside docker on `4317` and `4318`
  - Zipkin inside/outside docker on `9411`
- A Prometheus server which receives OpenTelemetry Collector metrics is available at http://localhost:9090.

## Labs

- Lab 101: Instrument an application. Complete one or more of the following:
  - [Instrument a Java application](java#lab-101-automatically-instrumenting-this-application-with-opentelemetry)
  - [Instrument a Python application](python#lab-101-automatically-instrumenting-this-application-with-opentelemetry)
  - [Instrument a NodeJS application](node#lab-101-automatically-instrumenting-this-application-with-opentelemetry)
  - [Instrument a Go application](go#lab-101-automatically-instrumenting-this-application-with-opentelemetry)
- Lab 102: Reconfigure applications to leverage the OpenTelemetry Collector.
  Complete one or more of the following:
  - [Reconfigure a Java application](java#lab-102-send-data-to-the-opentelemetry-collector)
  - [Reconfigure a Python application](python#lab-102-send-data-to-the-opentelemetry-collector)
  - [Reconfigure a NodeJS application](node#lab-102-send-data-to-the-opentelemetry-collector)
  - [Reconfigure a Go application](go#lab-102-send-data-to-the-opentelemetry-collector)
- Lab 103: Change context propagation format. Complete one (only one) of the
  following (requires you to instrument at least two applications):
  - [Change context propagation in a Java application](java#lab-103-instrument-another-application-and-call-it)
  - [Change context propagation in a Python application](python#lab-103-instrument-another-application-and-call-it)
  - [Change context propagation in a NodeJS application](node#lab-103-instrument-another-application-and-call-it)
  - [Change context propagation in a Go application](go#lab-103-instrument-another-application-and-call-it)
- [Lab 201: Processing trace data with the OpenTelemetry Collector](docker/lab200.md#lab-201-processing-trace-data-with-the-opentelemetry-collector)
- [Lab 202: Collecting and processing metric data with the OpenTelemetry Collector](docker/lab200.md#lab-202-collecting-and-processing-metric-data-with-the-opentelemetry-collector)
- [Lab 203: Collecting and processing log data with the OpenTelemetry Collector](docker/lab200.md#lab-203-collecting-and-processing-log-data-with-the-opentelemetry-collector)
