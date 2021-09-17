TODO

Notes:

to build/run in docker:

` $ ./gradlew jibDockerBuilder `


` $ docker run -p8080:8080 otel-workshop:1.0-SNAPSHOT `

OR better yet, for the whole suite:

` $ cd ../docker; docker-compose up `

To hit the url:

` $ curl http://localhost:8080/hello/proxy/java `

