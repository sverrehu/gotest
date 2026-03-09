// Example: gg + gogpu integration via ggcanvas
//
// This example demonstrates rendering 2D graphics with gg
// directly into a gogpu window using the ggcanvas integration package.
//
// Architecture:
//
//	gg.Context (draw) → ggcanvas.Canvas → gogpu.Context (GPU) → Window
//
// The example showcases all four GPU rendering tiers:
//
//	Tier 1 (SDF):           circles, rounded rectangles
//	Tier 2a (Convex):       triangle, pentagon, hexagon
//	Tier 2b (Stencil+Cover): star shape, curved paths
//	Tier 4 (MSDF text):     title text, FPS counter
//
// Rendering mode: event-driven with animation token.
// Uses ContinuousRender=false + StartAnimation() to render at VSync
// only while animation is active. Press Space to pause/resume.
//
// Requirements:
//   - gogpu v0.22.0+
//   - gg v0.31.0+
package main

import (
	"log"
	"time"

	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu" // Register GPU accelerator (SDF + MSAA 4x + MSDF text)
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
)

func main() {
	const width, height = 800, 600

	app := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle("GoGPU + gg: Four-Tier GPU Rendering").
		WithSize(width, height).
		WithContinuousRender(false)) // Event-driven: 0% CPU when paused

	var canvas *ggcanvas.Canvas
	var animToken *gogpu.AnimationToken
	var frame int
	paused := false
	startTime := time.Now()

	app.OnDraw(func(dc *gogpu.Context) {
		if frame == 0 {
			log.Printf("Backend: %s", dc.Backend())
			// Start animation — renders at VSync while token is alive.
			animToken = app.StartAnimation()
			log.Printf("Animation started (Space to pause/resume)")
		}

		w, h := dc.Width(), dc.Height()
		if w <= 0 || h <= 0 {
			return
		}

		// No dc.Clear() needed — gg renders directly to surface.
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
			log.Printf("Canvas created: %dx%d", w, h)
		}

		cw, ch := canvas.Size()
		if cw != w || ch != h {
			if err := canvas.Resize(w, h); err != nil {
				log.Printf("Resize error: %v", err)
			}
			cw, ch = w, h
		}

		elapsed := time.Since(startTime).Seconds()
		if err := canvas.Draw(func(cc *gg.Context) {
			renderFrame(cc, elapsed, cw, ch, frame)
		}); err != nil {
			log.Printf("Draw error: %v", err)
		}

		// Render directly to surface (zero-copy, no readback).
		sv := dc.SurfaceView()
		sw, sh := dc.SurfaceSize()
		if err := canvas.RenderDirect(sv, sw, sh); err != nil {
			log.Printf("Frame %d: RenderDirect error: %v", frame, err)
		}
		frame++
	})

	// Space toggles animation pause/resume — demonstrates three-state model.
	app.EventSource().OnKeyPress(func(key gpucontext.Key, _ gpucontext.Modifiers) {
		if key != gpucontext.KeySpace {
			return
		}
		paused = !paused
		if paused {
			if animToken != nil {
				animToken.Stop()
				animToken = nil
			}
			log.Printf("Paused (0%% CPU idle, press Space to resume)")
		} else {
			animToken = app.StartAnimation()
			log.Printf("Resumed")
		}
	})

	// GPU resources are automatically cleaned up on shutdown:
	// - ggcanvas.Canvas auto-registers with App's ResourceTracker
	// - App.Run() calls tracker.CloseAll() before Renderer.Destroy()
	// - OnClose is still available for additional cleanup (e.g., accelerator)
	app.OnClose(func() {
		if animToken != nil {
			animToken.Stop()
		}
		// Close accelerator: drains GPU queue and destroys session
		// resources (persistent buffers, textures) while the device is alive.
		gg.CloseAccelerator()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// renderFrame draws animated 2D graphics demonstrating all four GPU tiers.
func renderFrame(cc *gg.Context, elapsed float64, width, height int, frame int) {
	cc.SetRGBA(0, 0, 0, 0)
	cc.Clear()
	cc.SetRGB(1, 0, 1)
	renderShip(cc, 30, 30)
}

func renderShip(cc *gg.Context, x, y int) {
	const l = 30
	tipX := x + l/2
	tipY := y
	backLeftX := x - l/2
	backLeftY := y - l/2
	backRightX := backLeftX
	backRightY := y + l/2

	cc.MoveTo(float64(tipX), float64(tipY))
	cc.LineTo(float64(backLeftX), float64(backLeftY))
	cc.LineTo(float64(backRightX), float64(backRightY))
	cc.ClosePath()
	cc.Fill()
}
