# Python Workshop

## Getting Started

Start the suite of interconnected apps:

```bash
cd ../docker; docker-compose up --build
```

Then go to the url from outside the docker compose environment:

```bash
curl -i http://localhost:5000/hello/proxy/python
```

You can also just build and start this application (instead of the full suite):

```bash
docker build -t otel-workshop-python .
docker run -p5000:80 otel-workshop-python
```

## Lab 101: Automatically instrumenting this application with OpenTelemetry

The [OpenTelemetry getting started
documentation](https://github.com/open-telemetry/opentelemetry-python#getting-started)
covers the required steps.

In summary, your task is to:

- Install the required dependencies. Add the following to `requirements.txt`:
```bash
opentelemetry-distro==0.24b0
opentelemetry-exporter-jaeger-thrift==1.5.0
opentelemetry-instrumentation==0.24b0
```
- Add this line to the Dockerfile: `RUN opentelemetry-bootstrap -a=install` below the `RUN pip ...` command.
- Change the `CMD` line in the Dockerfile to: `CMD [ "opentelemetry-instrument", "python3", "app.py" ]`
- In the `docker/docker-compose.yml` file, add the following `environment` section to the python service:
```yaml
    environment:
      OTEL_TRACES_EXPORTER: jaeger_thrift
      OTEL_EXPORTER_JAEGER_ENDPOINT: http://jaeger:14268/api/traces
      OTEL_SERVICE_NAME: workshop-python-app
```
- Rebuild and restart docker-compose, as above.
- Run `curl -i http://localhost:5000/` -- what happened?
- Run `curl -i http://localhost:5000/hello` -- what happened?
- Run `curl -i http://localhost:5000/hello/proxy/node/python` -- what happened?

Assuming everything is configured properly, you should see the
`workshop-python-app` service in [Jaeger](http://localhost:16686).

## Lab 102: Send data to the OpenTelemetry Collector

By default, the OTLP gRPC exporter sends data to `localhost:4317`. Let's change
to this exporter to take advantage of the OpenTelemety Collector that is
running:

- In the `requirements.txt` file, add `opentelemetry-exporter-otlp-proto-grpc==1.5.0`
- In the `docker/docker-compose.yml` file, change the `environment` section of the python service:
```yaml
    environment:
      ...
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otelcol:4317
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:5000/` -- what happened?
- Run `curl -i http://localhost:5000/hello` -- what happened?
- Run `curl -i http://localhost:5000/hello/proxy/node/python` -- what happened?

> Question: Why is the exporter endpoint set to `otelcol` instead of using `localhost`?

Assuming everything is configured properly, you should see the
`workshop-python-app` service in [Jaeger](http://localhost:16686).

## Lab 103: Instrument another application and call it

Instrument another application (another language). Follow Lab 101 for that
other application. Then run `curl -i
http://localhost:5000/hello/proxy/<otherLanguage>`.

> Question: What do you see in Jaeger now?

Let's change the context propagation mechanism. In this application (not the
other) let's add one more environment variable:

- In the `docker/docker-compose.yml` file, change the `environment` section of the python service:
```yaml
    environment:
      ...
      OTEL_PROPAGATORS: b3multi
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:5000/hello/proxy/<otherLanguage>` -- what happened?

> Question: Check the `workshop-python-app` in Jaeger now. What happened? Why?
