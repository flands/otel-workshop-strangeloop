TODO

Notes:

to build/run in docker:

` $ ./gradlew build `

` $ docker run -p8080:80 otel-workshop `

OR better yet, for the whole suite (after doing a `./gradlew build`):

` $ cd ../docker; docker-compose up `

To hit the url from outside the container:

` $ curl http://localhost:8080/hello/proxy/java `

