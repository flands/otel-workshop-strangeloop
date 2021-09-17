package com.splunk.opentelemetry.workshop;

import spark.Request;
import spark.Route;
import spark.Spark;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

public class Main {

    public static final String HELLO_FROM_JAVA = "Hello from java\n";

    public static void main(String[] args) {
        HttpClient httpClient = HttpClient.newHttpClient();

        Main main = new Main();

        Spark.port(80);
        Spark.get("/hello", (request, response) -> main.handleHelloRequest());
        Spark.get("/hello/", (request, response) -> main.handleHelloRequest());
        Spark.get("/hello/proxy/*", (request, response) -> main.handleProxyRequest(httpClient, request));
    }

    private String handleHelloRequest() {
        return HELLO_FROM_JAVA;
    }

    private String handleProxyRequest(HttpClient httpClient, Request request) throws IOException, InterruptedException {
        String[] splat = request.splat();
        if (splat == null || splat.length == 0) {
            return HELLO_FROM_JAVA;
        }
        String proxyCommand = splat[0];
        String[] pieces = proxyCommand.split("/");

        StringBuilder rest =new StringBuilder();
        if (pieces.length > 0) {
            for (int i = 1; i< pieces.length; i++) {
                rest.append(pieces[i]).append("/");
            }
        }
        var uri = URI.create("http://" + pieces[0] + "/hello/proxy/" + rest);
        var httpResponse = httpClient.send(HttpRequest.newBuilder().GET().uri(uri).build(), HttpResponse.BodyHandlers.ofString());
        return httpResponse.body() + HELLO_FROM_JAVA;
    }
}
