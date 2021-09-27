# Go Workshop

## Getting Started

Start the suite of interconnected apps:

```bash
cd ../docker; docker-compose up --build
```

Then go to the url from outside the docker compose environment:

```bash
curl http://localhost:7070/hello/proxy/go
```

You can also just build and start this application (instead of the full suite):

```bash
docker build -t otel-workshop-go .
docker run -p7070:80 otel-workshop-go
```

## Lab 101: Automatically instrumenting this application with OpenTelemetry

Assuming everything is configured properly, you should see the
`workshop-go-app` service in [Jaeger](http://localhost:16686).

1. Similar to the [OpenTelemetry Go Jaeger
   example](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go),
   create a Jaeger exporter and register it with a new `TracerProvider`. That
   `TracerProvider` needs to be registered globally using
   [`otel.SetTracerProvider`](https://pkg.go.dev/go.opentelemetry.io/otel#SetTracerProvider).
   The Jaeger exporter does not need to be configured with any explicit URL,
   you will be setting that with environment variables.
2. With a configured and globally registered `TracerProvider`, use the [OpenTelemetry `net/http`
   instrumentation](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp).
   Similar to [the server
   example](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/net/http/otelhttp/example/server/server.go),
   wrap all the HTTP handlers with the `otelhttp.NewHandler`. This will trace
   all requests served by the server. Finally, update the HTTP request the
   proxy handler makes to used the traced
   [`otelhttp.Get`](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#Get)
   function.
3. In the [`docker/docker-compose.yml`](../docker/docker-compose.yml) file, add
   the following `environment` section to the go service:

    ```yaml
    environment:
      OTEL_EXPORTER_JAEGER_ENDPOINT: "http://jaeger:14268/api/traces"
      OTEL_SERVICE_NAME: workshop-go-app
    ```

4. Rebuild and restart docker-compose, as above.

You should now be able to run the same requests as before:

```sh
curl -i http://localhost:7070/hello
curl -i http://localhost:7070/hello/proxy/node/python/go/go
```

but now you should see the `workshop-go-app` service in [Jaeger](http://localhost:16686).

## Lab 102: Send data to the OpenTelemetry Collector

By default, the `otlptrace` gRPC exporter sends data on `localhost:4317`. Let's use
this exporter to take advantage of the OpenTelemety Collector that is running:

- Similar to the [OpenTelemetry Go OTel Collector
   example](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go),
   create an OTLP trace gRPC exporter and register it with a new `TracerProvider`. That
   `TracerProvider` needs to be registered globally using
   [`otel.SetTracerProvider`](https://pkg.go.dev/go.opentelemetry.io/otel#SetTracerProvider).
   The OTLP trace gRPC exporter does not need to be configured with any explicit URL,
   you will be setting that with environment variables.
- In the `docker/docker-compose.yml` file, change the `environment` section of the go service:
```yaml
    environment:
      ...
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otelcol:4317
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:7070/` -- what happened?
- Run `curl -i http://localhost:7070/hello` -- what happened?
- Run `curl -i http://localhost:7070/hello/proxy/node/python` -- what happened?

> Question: Why is the exporter endpoint set to `otelcol` instead of using `localhost`?

Assuming everything is configured properly, you should see the
`workshop-go-app` service in [Jaeger](http://localhost:16686).

## Lab 103: Instrument another application and call it

Instrument another application (another language). Follow Lab 101 for that
other application. Then run `curl -i
http://localhost:7070/hello/proxy/<otherLanguage>`.

> Question: What do you see in Jaeger now?

Let's change the context propagation mechanism. In this application (not the
other) let's add one more environment variable:

- In the `docker/docker-compose.yml` file, change the `environment` section of the go service:
```yaml
    environment:
      ...
      OTEL_PROPAGATORS: b3multi
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:7070/hello/proxy/<otherLanguage>` -- what happened?

> Question: Check the `workshop-go-app` in Jaeger now. What happened? Why?
