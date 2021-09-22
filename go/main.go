package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Paths the server listens on.
const (
	helloPath = "/hello"
	proxyPath = "/hello/proxy/"
)

// printHello writes the hello greeting to w.
func printHello(w http.ResponseWriter) {
	fmt.Fprintf(w, "Hello from Go\n")
}

// hello handles hello requests.
func hello(w http.ResponseWriter, req *http.Request) {
	printHello(w)
}

// proxyURL returns a URL string for the next service to proxy a request to.
//
// The orig URL is assumed to be a proxy request URL that contain services
// that need to be requested.
func proxyURL(orig *url.URL) string {
	u := &url.URL{Scheme: "http", Path: proxyPath}
	services := strings.TrimPrefix(orig.Path, proxyPath)
	switch parts := strings.SplitN(services, "/", 2); len(parts) {
	case 2:
		u.Path = u.Path + parts[1]
		fallthrough
	case 1:
		u.Host = parts[0]
	}
	return u.String()
}

// proxy handles proxy requests.
func proxy(w http.ResponseWriter, req *http.Request) {
	defer printHello(w)

	if req.URL.Path == proxyPath {
		return
	}

	resp, err := otelhttp.Get(req.Context(), proxyURL(req.URL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func initTracer() *trace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		log.Fatal(err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		append(
			resource.Default().Attributes(),
			semconv.ServiceNameKey.String("workshop-go-app"),
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

func main() {
	// Do not include timestamps in log output.
	log.SetFlags(0)
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	http.Handle(helloPath, otelhttp.NewHandler(http.HandlerFunc(hello), "/hello"))
	http.Handle(proxyPath, otelhttp.NewHandler(http.HandlerFunc(proxy), "/hello/proxy/*"))
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
