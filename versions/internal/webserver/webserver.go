package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	maven "github.com/sverrehu/gotest/versions/internal/repos"
)

type handler struct {
	target  string
	handler http.Handler
}

var handlers []handler

type mavenHandler struct {
}

func (h *mavenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pkg := r.PathValue("package")
	parts := regexp.MustCompile("[:/]").Split(pkg, -1)
	if len(parts) != 2 {
		sendBadRequest(w, "expected two parts, separated by ':' or '/' in maven package", pkg)
		return
	}
	releases, err := maven.GetReleases(parts[0], parts[1])
	if err != nil {
		sendInternalServerError(w, err, r.URL)
		return
	}
	jsonReleases, err := json.Marshal(releases)
	if err != nil {
		sendInternalServerError(w, err, r.URL)
		return
	}
	w.Write(jsonReleases)
}

func sendInternalServerError(w http.ResponseWriter, err error, url *url.URL) {
	log.Printf("internal server error for url: %v: %v", url, err.Error())
	w.WriteHeader(http.StatusBadRequest)
}

func sendBadRequest(w http.ResponseWriter, message string, pkg string) {
	log.Printf("bad request: %s, got: %s", message, pkg)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
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
		{target: "/maven", handler: &mavenHandler{}},
	}
}
