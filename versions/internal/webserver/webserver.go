package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sverrehu/gotest/versions/internal/repos"
)

type handler struct {
	target  string
	handler http.Handler
}

var handlers []handler

type commonReleasesHandler struct {
	h repos.ReleasesFetcher
}

func (h *commonReleasesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pkg := r.PathValue("package")
	releases, err := h.h.GetReleases(pkg)
	if err != nil {
		var re *repos.ReleasesFetcherError
		ok := errors.As(err, &re)
		if ok && re.IsParameterError {
			sendBadRequest(w, err.Error(), r.URL)
		} else {
			sendInternalServerError(w, err, r.URL)
		}
		return
	}
	jsonReleases, err := json.Marshal(releases)
	if err != nil {
		sendInternalServerError(w, err, r.URL)
		return
	}
	_, err = w.Write(jsonReleases)
	if err != nil {
		log.Printf("error writing response for url: %v: %v", r.URL, err.Error())
	}
}

func sendInternalServerError(w http.ResponseWriter, err error, url *url.URL) {
	log.Printf("internal server error for url: %v: %v", url, err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func sendBadRequest(w http.ResponseWriter, message string, url *url.URL) {
	log.Printf("bad request: %s", message)
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
	if err != nil {
		sendInternalServerError(w, err, url)
	}
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
		{target: "/maven", handler: &commonReleasesHandler{&repos.MavenReleasesFetcher{}}},
		{target: "/dockerhub", handler: &commonReleasesHandler{&repos.OCIReleasesFetcher{}}},
		{target: "/github-releases", handler: &commonReleasesHandler{&repos.GitHubReleasesFetcher{}}},
	}
}
