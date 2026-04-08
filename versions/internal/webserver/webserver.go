package webserver

import (
	"fmt"
	"net/http"
	"strconv"
)

type handler struct {
	target  string
	handler http.Handler
}

var handlers []handler

type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%s", r.PathValue("package"))))
}

func Run() error {
	port := 8086
	mux := http.NewServeMux()
	for _, h := range handlers {
		fmt.Printf("Adding handler for %s\n", h.target)
		mux.Handle(h.target+"/{package...}", h.handler)
	}
	fmt.Printf("Starting server at port %d. Ctrl-C to abort.\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	return err
}

func init() {
	handlers = []handler{
		{target: "/maven", handler: &indexHandler{}},
	}
}
