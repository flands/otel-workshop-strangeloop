package com.splunk.opentelemetry.workshop;

import spark.Request;
import spark.Response;
import spark.Spark;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Arrays;
import java.util.List;

public class Main {

    public static final String HELLO_FROM_JAVA = "Hello from java\n";

    private final HttpClient httpClient;

    public Main(HttpClient httpClient) {
        this.httpClient = httpClient;
    }

    public static void main(String[] args) {
        HttpClient httpClient = HttpClient.newHttpClient();

        Main main = new Main(httpClient);

        Spark.port(80);
        Spark.get("/hello", (request, response) -> main.handleHelloRequest());
        Spark.get("/hello/proxy/*", (request, response) -> main.handleProxyRequest(request, response));
    }

    private String handleHelloRequest() {
        return HELLO_FROM_JAVA;
    }

    private String handleProxyRequest(Request request, Response response) throws IOException, InterruptedException {
        String[] splat = request.splat();
        if (splat == null || splat.length == 0) {
            response.status(404);
            return "NOT FOUND";
        }
        String proxyCommand = splat[0];

        List<String> pieces = Arrays.asList(proxyCommand.split("/"));
        List<String> rest = pieces.subList(1, pieces.size());
        String proxyTarget = pieces.get(0);
        URI uri;
        if (rest.size() == 1) {
            uri = URI.create("http://" + proxyTarget + "/hello");
        } else {
            uri = URI.create("http://" + proxyTarget + "/hello/proxy/" + String.join("/", rest));
        }
        var httpResponse = httpClient.send(HttpRequest.newBuilder().GET().uri(uri).build(), HttpResponse.BodyHandlers.ofString());
        return httpResponse.body() + HELLO_FROM_JAVA;
    }
}
