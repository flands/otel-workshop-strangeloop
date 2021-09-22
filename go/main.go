package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func printHello(w http.ResponseWriter) {
	fmt.Fprintf(w, "Hello from Go\n")
}

func hello(w http.ResponseWriter, req *http.Request) {
	printHello(w)
}

func proxyURL(orig *url.URL) *url.URL {
	path := strings.TrimSuffix(orig.Path, "/")
	idx := strings.LastIndex(path, "/")
	return &url.URL{Scheme: "http", Host: path[idx+1:], Path: path[:idx]}
}

func proxy(w http.ResponseWriter, req *http.Request) {
	defer printHello(w)

	u := proxyURL(req.URL)
	if u.Host == "proxy" {
		// Handle the request to /hello/proxy/ directly."
		return
	}
	resp, err := http.Get(u.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/hello/proxy", hello)
	http.HandleFunc("/hello/proxy/", proxy)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
