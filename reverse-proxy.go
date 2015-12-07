package main

import (
	"net/http"
	"log"
	"fmt"
	"net/url"
	"net/http/httputil"
	"time"
	"os"
	"flag"
	"errors"
)

func NewSingleHostReverseProxy(target *url.URL) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		targetQuery := target.RawQuery
		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			//req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		}
		(&httputil.ReverseProxy{
			Director: director,
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}

func main() {
	port := flag.String("port", "", "public server port, e.g. ':80'")
	target := flag.String("target", "", "default target url, e.g. 'http://127.0.0.1:8080'")

	flag.Parse()

	if (*port == "" || *target == "") {
		if (!(*port == "" && *target == "")) {
			if (*port == "") {
				fmt.Fprintf(os.Stderr, "error: %v\n", errors.New("please provide a 'port'"))
			}
			if (*target == "") {
				fmt.Fprintf(os.Stderr, "error: %v\n", errors.New("please provide the 'target'"))
			}
		}

		flag.Usage()
		os.Exit(1)
	}

	url, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("couldn't parse target url '%v'", *target)
		flag.Usage()
		os.Exit(1)
	}

	http.HandleFunc("/", NewSingleHostReverseProxy(url))

	log.Printf("ready to serve on %v", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
