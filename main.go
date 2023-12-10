package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseProxy(target *url.URL) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	targetURL := "http://localhost:8000"
	target, _ := url.Parse(targetURL)

	proxyHandler := reverseProxy(target)

	server := http.Server{
		Addr:    ":8080",
		Handler: proxyHandler,
	}

	fmt.Println("Starting proxy on :8080")
	fmt.Println("Forwarding to:", targetURL)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Failed on create server:", err)
	}
}
