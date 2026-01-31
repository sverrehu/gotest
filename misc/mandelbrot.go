package main

import (
	"math/big"

	"thathost.com/golang/gotest/misc/mandelbrot"
)

func main() {
	WIDTH := 3840
	HEIGHT := 2160
	coordCenterX := big.NewFloat(-0.75)
	coordCenterY := big.NewFloat(0.0)
	coordHeight := big.NewFloat(2.0)
	maxIterations := int16(10000)
	result := mandelbrot.Calculate(int32(WIDTH), int32(HEIGHT), coordCenterX, coordCenterY, coordHeight, maxIterations, false)
	err := mandelbrot.Save(result, "/tmp/latest.mandel")
	if err != nil {
		panic(err)
	}
}
