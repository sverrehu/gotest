package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Go webserver!"))
}

func main() {
	port := 8086
	mux := http.NewServeMux()
	mux.Handle("/", &indexHandler{})
	fmt.Printf("Starting server at port %d. Ctrl-C to abort.\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		panic(err)
	}
}
