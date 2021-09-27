# Node Workshop

## Getting Started

Start the suite of interconnected apps:

```bash
cd ../docker; docker-compose up --build
```

Then go to the url from outside the docker compose environment:

```bash
curl -i http://localhost:3000/hello/proxy/node
```

You can also just build and start this application (instead of the full suite):

```bash
docker build -t otel-workshop-node .
docker run -p3000:80 otel-workshop-node
```

## Lab 101: Automatically instrumenting this application with OpenTelemetry

The [OpenTelemetry getting started
documentation](https://github.com/open-telemetry/opentelemetry-js/blob/main/getting-started/README.md#trace-your-nodejs-application)
covers the required steps.

In summary, your task is to:

- Install the required OpenTelemetry libraries
```bash
npm install \
  @opentelemetry/api \
  @opentelemetry/sdk-trace-base \
  @opentelemetry/sdk-trace-node \
  @opentelemetry/instrumentation-http \
  @opentelemetry/instrumentation-express \
  @opentelemetry/instrumentation-grpc \
  @opentelemetry/exporter-jaeger
```

If you don't have `npm` installed locally, you can directly add the following dependencies to the `package.json` file,
under the "node-fetch" entry:
```
    "@opentelemetry/api": "^1.0.3",
    "@opentelemetry/exporter-jaeger": "^0.25.0",
    "@opentelemetry/instrumentation-express": "^0.25.0",
    "@opentelemetry/instrumentation-grpc": "^0.25.0",
    "@opentelemetry/instrumentation-http": "^0.25.0",
    "@opentelemetry/sdk-trace-base": "^0.25.0",
    "@opentelemetry/sdk-trace-node": "^0.25.0"
```

- Initialize a global trace. Create a file named `tracing.js` and add the following code:
```javascript
'use strict'

const process = require('process');
const opentelemetry = require('@opentelemetry/sdk-node');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-base');
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
const { diag, DiagConsoleLogger, DiagLogLevel} = require("@opentelemetry/api");

diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG)

console.log(process.env)
// configure the SDK to export telemetry data to the console
// enable all auto-instrumentations from the meta package
const traceExporter = new JaegerExporter();
const sdk = new opentelemetry.NodeSDK({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'workshop-node-app',
    }),
    traceExporter,
    instrumentations: [getNodeAutoInstrumentations()]
});

// initialize the SDK and register with the OpenTelemetry API
// this enables the API to record telemetry
sdk.start()
    .then(() => console.log('Tracing initialized'))
    .catch((error) => console.log('Error initializing tracing', error));

// gracefully shut down the SDK on process exit
process.on('SIGTERM', () => {
    sdk.shutdown()
        .then(() => console.log('Tracing terminated'))
        .catch((error) => console.log('Error terminating tracing', error))
        .finally(() => process.exit(0));
});
```
- In the `Dockerfile`:
  - Add a line to add the new `tracing.js` to the docker build :
```dockerfile
ADD tracing.js /app/
```
  - change the `CMD` line to read:
```dockerfile
CMD ["node", "-r", "./tracing.js", "app.js"]
```

- In the `docker/docker-compose.yml` file, add the following `environment` section to the node service:
```yaml
    environment:
      OTEL_EXPORTER_JAEGER_ENDPOINT: "http://jaeger:14268/api/traces"
      OTEL_SERVICE_NAME: "workshop-node-app"
```
- Rebuild and restart docker-compose, as above.
- Run `curl -i http://localhost:3000/` -- what happened?
- Run `curl -i http://localhost:3000/hello` -- what happened?
- Run `curl -i http://localhost:3000/hello/proxy/node/python` -- what happened?

Assuming everything is configured properly, you should see the
`workshop-node-app` service in [Jaeger](http://localhost:16686).

## Lab 102: Send data to the OpenTelemetry Collector

By default, the OTLP gRPC exporter sends data to `localhost:4317`. Let's use
this exporter to take advantage of the OpenTelemety Collector that is running:

- Install the required dependency
```bash
npm install --save @opentelemetry/exporter-otlp-grpc
```
- Update tracing.js to contain:
```bash
const { BasicTracerProvider, SimpleSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { OTLPTraceExporter } =  require('@opentelemetry/exporter-otlp-grpc');

const provider = new BasicTracerProvider();
const exporter = new OTLPTraceExporter();
provider.addSpanProcessor(new SimpleSpanProcessor(exporter));
```
- In the `docker/docker-compose.yml` file, change the `environment` section of the node service:
```yaml
    environment:
      ...
      OTEL_EXPORTER_OTLP_ENDPOINT: http://otelcol:4317
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:3000/` -- what happened?
- Run `curl -i http://localhost:3000/hello` -- what happened?
- Run `curl -i http://localhost:3000/hello/proxy/node/python` -- what happened?

> Question: Why is the exporter endpoint set to `otelcol` instead of using `localhost`?

Assuming everything is configured properly, you should see the
`workshop-node-app` service in [Jaeger](http://localhost:16686).

## Lab 103: Instrument another application and call it

Instrument another application (another language). Follow Lab 101 for that
other application. Then run `curl -i
http://localhost:3000/hello/proxy/<otherLanguage>`.

> Question: What do you see in Jaeger now?

Let's change the context propagation mechanism. In this application (not the
other) let's add one more environment variable:

- In the `docker/docker-compose.yml` file, change the `environment` section of the node service:
```yaml
    environment:
      ...
      OTEL_PROPAGATORS: b3multi
```
- Restart docker-compose as above.
- Run `curl -i http://localhost:3000/hello/proxy/<otherLanguage>` -- what happened?

> Question: Check the `workshop-node-app` in Jaeger now. What happened? Why?
