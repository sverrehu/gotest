package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go webserver!"))
}

func main() {
	reg := prometheus.NewRegistry()
	port := 8086
	mux := http.NewServeMux()
	mux.Handle("/", &indexHandler{})
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	fmt.Printf("Starting server at port %d. Ctrl-C to abort.\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		panic(err)
	}
}
