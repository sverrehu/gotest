package main

import (
	"log"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
)

const quadrantWidth, quadrantHeight = 1500, 850

func main() {
	app := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle("SHH Space Game").
		WithSize(quadrantWidth, quadrantHeight).
		WithContinuousRender(false))

	var canvas *ggcanvas.Canvas
	var animToken *gogpu.AnimationToken
	var frame int

	app.OnDraw(func(dc *gogpu.Context) {
		if frame == 0 {
			animToken = app.StartAnimation()
		}
		w, h := dc.Width(), dc.Height()
		if w <= 0 || h <= 0 {
			return
		}
		if canvas == nil {
			provider := app.GPUContextProvider()
			if provider == nil {
				return
			}
			var err error
			canvas, err = ggcanvas.New(provider, w, h)
			if err != nil {
				log.Fatalf("Failed to create canvas: %v", err)
			}
		}
		cw, ch := canvas.Size()
		if cw != w || ch != h {
			if err := canvas.Resize(w, h); err != nil {
				log.Printf("Resize error: %v", err)
			}
			cw, ch = w, h
		}
		err := canvas.Draw(func(cc *gg.Context) {
			quadrantRenderFrame(cc)
		})
		if err != nil {
			log.Printf("Draw error: %v", err)
		}
		err = canvas.Render(dc.RenderTarget())
		if err != nil {
			log.Printf("Frame %d: Render error: %v", frame, err)
		}
		app.RequestRedraw()
		frame++
	})

	app.OnClose(func() {
		if animToken != nil {
			animToken.Stop()
		}
		gg.CloseAccelerator()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

var ballX = quadrantWidth / 2.0
var ballY = quadrantHeight / 2.0
var ballDx = -3.0
var ballDy = -2.0

const ballRadius = 10

func quadrantRenderFrame(cc *gg.Context) {
	cc.ClearWithColor(gg.RGBA2(0, 0, 0, 1))
	_ = cc.Fill()
	cc.SetRGB(1, 1, 0.5)
	cc.DrawCircle(ballX, ballY, ballRadius)
	_ = cc.Fill()
	ballX += ballDx
	if ballX >= quadrantWidth-ballRadius || ballX < ballRadius {
		ballX -= ballDx
		ballDx = -ballDx
	}
	ballY += ballDy
	if ballY >= quadrantHeight-ballRadius || ballY < ballRadius {
		ballY -= ballDy
		ballDy = -ballDy
	}
}
