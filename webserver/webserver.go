package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type indexHandler struct{}

type metrics struct {
	numRequests prometheus.Counter
}

var m *metrics

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		numRequests: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "webserver_requests_total",
			Help: "Total number of requests received by the webserver",
		}),
	}
	return m
}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go webserver!\n"))
	m.numRequests.Inc()
}

func main() {
	reg := prometheus.NewRegistry()
	m = newMetrics(reg)
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
