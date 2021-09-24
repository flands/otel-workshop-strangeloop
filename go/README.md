# Go Workshop

## Getting Started

Start the suite of interconnected apps (it is OK to just instrument this service):

```sh
cd ../docker; docker-compose up --build
```

Then go to the url from outside the docker compose environment:

```sh
curl http://localhost:7070/hello/proxy/go
```

You can also just build and start this application (instead of the full suite):

```sh
docker build -t otel-workshop-go .
docker run -p7070:80 otel-workshop-go
```

## Lab 101: Instrumenting this Application with OpenTelemetry

1. Similar to the [OpenTelemetry Jaeger
   example](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go),
   create a Jaeger exporter and register it with a new `TracerProvider`. That
   `TracerProvider` needs to be registered globally using
   [`otel.SetTracerProvider`](https://pkg.go.dev/go.opentelemetry.io/otel#SetTracerProvider).
   The Jaeger exporter does not need to be configured with any explicit URL,
   you will be setting that with environment variables.
2. With a configured and globally registered `TracerProvider`, now you need to
   use the [OpenTelemetry `net/http`
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

You should now be able to run the same requests as before, e.g.

```sh
curl http://localhost:7070/hello
curl http://localhost:7070/hello/proxy/node/python/java/go
```

but now you should see the `workshop-go-app` service in Jaeger.
