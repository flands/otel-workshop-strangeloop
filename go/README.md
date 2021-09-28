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

## Lab 101: Instrumenting this application with OpenTelemetry

Assuming everything is configured properly, you should see the
`workshop-go-app` service in [Jaeger](http://localhost:16686).

1. Add the following to the `go.mod` file:
```
require (
		go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.24.0
		go.opentelemetry.io/otel v1.0.0
		go.opentelemetry.io/otel/exporters/jaeger v1.0.0
		go.opentelemetry.io/otel/sdk v1.0.0
		go.opentelemetry.io/otel/trace v1.0.0 // indirect
)
```
2. In `main.go`, add the following imports:
```
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
```
3. Create a new function to initialize the OpenTelemetry Tracing SDK:
```
func initTracer() *trace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		log.Fatal(err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		append(
			resource.Default().Attributes(),
		)...,
	)

	tp := trace.NewTracerProvider(
		// Sample all traces for this demo application.
		// (not recommended for production).
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp
}
```
4. In the `main()` function, add a call to the method you just created, after the `log.SetFlags(0)` call:
```
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
```
5. At this point, the code is ready to have the http server and client instrumentation added to it.
   - Use the `otelhttp` implementation to instrument the http client call. in the `proxy` function, replace the call to `http.Get(...)` with this:
   ```
       resp, err := otelhttp.Get(req.Context(), proxyURL(req.URL))
   ```
   - In the `main()` function, replace the registration of the http handlers with the following:
   ```
       http.Handle(helloPath, otelhttp.NewHandler(http.HandlerFunc(hello), "/hello"))
       http.Handle(proxyPath, otelhttp.NewHandler(http.HandlerFunc(proxy), "/hello/proxy/*"))
   ```
6. In the [`docker/docker-compose.yml`](../docker/docker-compose.yml) file, add
   the following `environment` section to the go service:

    ```yaml
    environment:
      OTEL_EXPORTER_JAEGER_ENDPOINT: "http://jaeger:14268/api/traces"
      OTEL_SERVICE_NAME: workshop-go-app
    ```
7. Rebuild and restart docker-compose, as above.

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
