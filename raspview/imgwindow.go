package main

import (
	"image"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/integration/ggcanvas"
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
	w.app = gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle("Image Viewer").
		WithSize(800, 600))
	var canvas *ggcanvas.Canvas
	w.app.OnDraw(func(dc *gogpu.Context) {
		width, height := dc.Width(), dc.Height()
		if canvas == nil || canvas.Width() != width || canvas.Height() != height {
			canvas, _ = ggcanvas.New(w.app.GPUContextProvider(), width, height)
		}
		ctx := canvas.Context()
		ctx.ClearWithColor(gg.Black)
		if w.img != nil {
			imgBuf := gg.ImageBufFromImage(*w.img)
			ctx.DrawImage(imgBuf, float64(width-imgBuf.Width())/2.0, float64(height-imgBuf.Height())/2.0)
		}
		canvas.RenderTo(dc.AsTextureDrawer())
	})
	w.app.Run()
}
