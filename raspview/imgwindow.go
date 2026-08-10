package main

import (
	"image"

	"github.com/gogpu/gogpu"
)

type ImageWindow struct {
	app *gogpu.App
	img *image.Image
}

func NewImageWindow() *ImageWindow {
	return &ImageWindow{}
}

func (w *ImageWindow) SetImage(img *image.Image) {
	w.img = img
	w.openOrRedraw()
}

func (w *ImageWindow) openOrRedraw() {
	if w.app == nil {
		w.open()
	}
}

func (w *ImageWindow) open() {

}
