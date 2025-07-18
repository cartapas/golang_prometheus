package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	// HTTP endpoint using the promhttp library
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
