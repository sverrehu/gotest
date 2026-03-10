package main

import (
	"github.com/gogpu/gg"
)

func main() {
	dc := gg.NewContext(800, 600)
	defer dc.Close()

	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 0, 1)
	renderShip2(dc, 30, 30)

	err := dc.SavePNG("output.png")
	if err != nil {
		panic(err)
	}
}

func renderShip2(cc *gg.Context, x, y int) {
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
