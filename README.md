# OpenTelemetry StrangeLoop Workshop

Prerequisites:
- Java 11
  - Note: this could probably be mitigated if we build the app in docker itself, rather than externally.
- Docker w/ docker-compose

Notes on getting the workshop up & running:

```
$ cd java; ./gradlew build
$ cd ../docker
$ docker-compose up
```

Exercise the API:

```
$ curl -i localhost:3000/hello
$ curl -i localhost:8080/hello/proxy/java/node/java
```

View Traces in Jaeger : http://localhost:16686
