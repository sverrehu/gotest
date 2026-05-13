package main

import (
	"log"
	"os"

	"github.com/gogpu/gg"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gg/text"
	"github.com/gogpu/gogpu"
)

const width, height = 300, 300

var font text.Face

func StartUI() {
	app := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle("").
		WithSize(width, height).
		WithContinuousRender(false))

	var canvas *ggcanvas.Canvas
	var animToken *gogpu.AnimationToken
	var frame int

	fontSource := loadFontSource()
	font = fontSource.Face(48)

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
		if err := canvas.Draw(func(cc *gg.Context) {
			renderFrame(cc)
		}); err != nil {
			log.Printf("Draw error: %v", err)
		}
		if err := canvas.Render(dc.RenderTarget()); err != nil {
			log.Printf("Frame %d: Render error: %v", frame, err)
		}
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

func renderFrame(cc *gg.Context) {
	cc.SetRGB(0, 0, 0)
	cc.Clear()
	cc.SetRGB(1, 0, 1)
	cc.MoveTo(width/2, 0)
	cc.LineTo(width/2, height-1)
	cc.Stroke()
	cc.MoveTo(0, height/2)
	cc.LineTo(width-1, height/2)
	cc.Stroke()
	cc.SetRGB(1, 1, 1)
	cc.SetFont(font)
	cc.DrawStringAnchored("foobar", width/2, height/2, 0.5, 0.5)

}

func loadFontSource() *text.FontSource {
	fontPath := findSystemFont()
	if fontPath == "" {
		log.Fatal("No system font found.")
		return nil
	}
	source, err := text.NewFontSourceFromFile(fontPath)
	if err != nil {
		log.Fatalf("Failed to load font %s: %v", fontPath, err)
		return nil
	}
	log.Printf("Loaded font: %s", source.Name())
	return source
}

func findSystemFont() string {
	candidates := []string{
		// Linux
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/usr/share/fonts/liberation/LiberationSans-Regular.ttf",
		// macOS
		"/Library/Fonts/Arial.ttf",
		"/System/Library/Fonts/Supplemental/Arial.ttf",
		"/System/Library/Fonts/Supplemental/Courier New.ttf",
		"/System/Library/Fonts/Monaco.ttf",
		// Windows
		"C:\\Windows\\Fonts\\arial.ttf",
		"C:\\Windows\\Fonts\\calibri.ttf",
		"C:\\Windows\\Fonts\\segoeui.ttf",
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func main() {
	StartUI()
}
