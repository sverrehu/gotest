package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"net/http"
	"slices"
	"strconv"
)

// curl -X POST -H "Content-type: image/jpeg" --data-binary @"$HOME/Pictures/statements/BushIsrael.png" http://localhost:8086/img

type indexHandler struct{}

type imageHandler struct{}

var imgWindow *ImageWindow

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("POST a binary image to /img\n"))
}

func (h *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(content)
	if !supportedFormat(contentType) {
		http.Error(w, "Unsupported content: "+contentType, http.StatusBadRequest)
		return
	}
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		http.Error(w, "Unable to decode image", http.StatusBadRequest)
		return
	}
	showImage(&img)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK\n"))
}

func supportedFormat(f string) bool {
	supported := []string{"image/jpeg", "image/png"}
	return slices.Contains(supported, f)
}

func showImage(img *image.Image) {
	imgWindow.SetImage(img)
}

func startWebserver() {
	port := 8086
	mux := http.NewServeMux()
	mux.Handle("/", &indexHandler{})
	mux.Handle("POST /img", &imageHandler{})
	fmt.Printf("Starting server at port %d. Ctrl-C to abort.\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		panic(err)
	}
}

func startImageViewer() {
	imgWindow = NewImageWindow()
}

func main() {
	startImageViewer()
	startWebserver()
}
