package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	resp, err := http.Get(proxyURL(req.URL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Do not include timestamps in log output.
	log.SetFlags(0)

	http.HandleFunc(helloPath, hello)
	http.HandleFunc(proxyPath, proxy)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
