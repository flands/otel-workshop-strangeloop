Notes:

to build/run in docker:

` $ docker build -t otel-workshop-java . `

` $ docker run -p8080:80 otel-workshop-java `

Or better yet, for the whole suite of interconnected apps:

` $ cd ../docker; docker-compose up --build `

To hit the url from outside the docker compose environment:

` $ curl http://localhost:8080/hello/proxy/java `

### How to instrument with OpenTelemetry's auto-instrumentation javaagent:

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

[otel-release]: https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest
[otel-latest-jar]: https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent-all.jar