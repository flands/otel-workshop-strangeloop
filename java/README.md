# Java Workshop

## Getting Started

Start the suite of interconnected apps:

```
$ cd ../docker; docker-compose up --build
```

Then go to the url from outside the docker compose environment:

```
$ curl http://localhost:8080/hello/proxy/java
```

You can also just build and start this application (instead of the full suite):

```
$ docker build -t otel-workshop-java .
$ docker run -p8080:80 otel-workshop-java
```

## Lab 101: Automatically instrumenting this application with OpenTelemetry

The [OpenTelemetry getting started
documentation](https://github.com/open-telemetry/opentelemetry-java-instrumentation#getting-started)
covers the required steps.

In summary, your task is to:

- Download the [opentelemetry-javaagent-all.jar][otel-latest-jar] from the [opentelemetry repository][otel-release]
- Place the downloaded agent jar in the `java` directory (where this README.md file is located).
- Add this line to the Dockerfile: `RUN cp opentelemetry-javaagent-all.jar /app/` below the other `RUN cp ...` command.
- In the `build.gradle` file, inside the `application` block, change the `applicationDefaultJvmArgs` to be
  `['-Xmx128M', '-javaagent:opentelemetry-javaagent-all.jar']`
- In the `docker/docker-compose.yml` file, add the following `environment` section to the java service:
```
    environment:
      OTEL_TRACES_EXPORTER: jaeger
      OTEL_EXPORTER_JAEGER_ENDPOINT: http://jaeger:14250
      OTEL_SERVICE_NAME: workshop-java-app
```
- Rebuild and restart docker-compose, as above.
- Run `curl -i http://localhost:8080/` -- what happened?
- Run `curl -i http://localhost:8080/hello` -- what happened?
- Run `curl -i http://localhost:8080/hello/proxy/node/python` -- what happened?

Assuming everything is configured properly, you should see the
`workshop-java-app` service in [Jaeger](http://localhost:16686).

## Lab 102: Send data to the OpenTelemetry Collector

By default, auto-instrumentation sends data via OTLP on `localhost:4317`. In
Lab 101 we overrode this via environment variables. Let's modify them to take
advantage of the OpenTelemety Collector that is running:

- In the `docker/docker-compose.yml` file, change the `environment` section of the java service:
```
    environment:
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otelcol:4317
      OTEL_SERVICE_NAME: workshop-java-app
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:8080/` -- what happened?
- Run `curl -i http://localhost:8080/hello` -- what happened?
- Run `curl -i http://localhost:8080/hello/proxy/node/python` -- what happened?

> Question: Why is the exporter endpoint set to `otelcol` instead of using `localhost`?

Assuming everything is configured properly, you should see the
`workshop-java-app` service in [Jaeger](http://localhost:16686).

## Lab 103: Instrument another application and call it

Instrument another application (another language). Follow Lab 101 for that
other application. Then run `curl -i
http://localhost:8080/hello/proxy/<otherLanguage>`.

> Question: What do you see in Jaeger now?

Let's change the context propagation mechanism. In this application (not the
other) let's add one more environment variable:

- In the `docker/docker-compose.yml` file, change the `environment` section of the java service:
```
    environment:
      ...
      OTEL_PROPAGATORS: b3multi
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:8080/hello/proxy/<otherLanguage>` -- what happened?

> Question: Check the `workshop-java-app` in Jaeger now. What happened? Why?

[otel-release]: https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest
[otel-latest-jar]: https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent-all.jar
