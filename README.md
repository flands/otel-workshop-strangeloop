# OpenTelemetry StrangeLoop Workshop

Prerequisites:
- Docker w/ docker-compose

Notes on getting the workshop up & running:

```
$ cd ../docker
$ docker-compose up --build
```

Java is listening outside docker on `8080`.
Python is listening outside docker on `5000`.
Node is listening outside docker on `3000`.
Go is listening outside docker on `7070`.

Exercise the API:

```
$ curl -i localhost:3000/hello
$ curl -i localhost:8080/hello/proxy/python/node/java
```

Once instrumented for OpenTelemetry, view traces in Jaeger : http://localhost:16686
