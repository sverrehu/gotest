package mandelbrot

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	immutableMath "thathost.com/golang/gotest/misc/internal"
)

func Calculate(resultWidth int32, resultHeight int32, coordCenterX *big.Float, coordCenterY *big.Float, coordHeight *big.Float, maxIterations int16, multiThreaded bool) *MandelbrotResult {
	two := big.NewFloat(2)
	coordWidth := getCoordWidth(resultWidth, resultHeight, coordHeight)
	coordWidthHalf := immutableMath.Divide(coordWidth, two)
	coordHeightHalf := immutableMath.Divide(coordHeight, two)
	coordMinX := immutableMath.Subtract(coordCenterX, coordWidthHalf)
	coordMaxX := immutableMath.Add(coordCenterX, coordWidthHalf)
	coordMinY := immutableMath.Subtract(coordCenterY, coordHeightHalf)
	coordMaxY := immutableMath.Add(coordCenterY, coordHeightHalf)
	return calculate(resultWidth, resultHeight, coordWidth, coordHeight, coordMinX, coordMaxX, coordMinY, coordMaxY, maxIterations, multiThreaded)
}

func getCoordWidth(pixelWidth int32, pixelHeight int32, coordHeight *big.Float) *big.Float {
	return immutableMath.Divide(immutableMath.Multiply(big.NewFloat(float64(pixelWidth)), coordHeight), big.NewFloat(float64(pixelHeight)))
}

func calculate(resultWidth int32, resultHeight int32, coordWidth *big.Float, coordHeight *big.Float, coordMinX *big.Float, coordMaxX *big.Float, coordMinY *big.Float, coordMaxY *big.Float, maxIterations int16, multiThreaded bool) *MandelbrotResult {
	result := NewMandelbrotResult(resultWidth, resultHeight, CalculationTypeDouble)
	result.CoordMinX.Copy(coordMinX)
	result.CoordMaxX.Copy(coordMaxX)
	result.CoordMinY.Copy(coordMinY)
	result.CoordMaxY.Copy(coordMaxY)
	result.MaxIterations = maxIterations
	start := time.Now()
	var waitGroup sync.WaitGroup
	for px, _ := range result.IterationCounts {
		bigPx := big.NewFloat(float64(px))
		bigResultWidth := big.NewFloat(float64(resultWidth))
		x0 := immutableMath.Add(coordMinX, immutableMath.Multiply(bigPx, immutableMath.Divide(coordWidth, bigResultWidth)))
		column := result.IterationCounts[px]
		if multiThreaded {
			waitGroup.Go(func() { calculateColumn(column, resultHeight, coordHeight, x0, coordMinY, maxIterations, result) })
		} else {
			calculateColumn(column, resultHeight, coordHeight, x0, coordMinY, maxIterations, result)
		}
	}
	if multiThreaded {
		waitGroup.Wait()
	}
	elapsed := time.Since(start)
	result.CalculationTimeMs = elapsed.Milliseconds()
	fmt.Printf("Num pixels escaped:  %d\n", result.NumEscaped)
	fmt.Printf("Num pixels infinite: %d\n", result.NumInfinite)
	fmt.Printf("Done in %g sec\n", float64(elapsed.Milliseconds()/1000.0)) // Single thread: 22.96, multiple threads: 3.6
	return result
}

func calculateColumn(column []int16, resultHeight int32, coordHeight *big.Float, x0 *big.Float, yMin *big.Float, maxIterations int16, result *MandelbrotResult) {
	for py, _ := range column {
		bigPy := big.NewFloat(float64(py))
		bigResultHeight := big.NewFloat(float64(resultHeight))
		y0 := immutableMath.Add(yMin, immutableMath.Multiply(bigPy, immutableMath.Divide(coordHeight, bigResultHeight)))
		column[int(resultHeight)-py-1] = calculateIterations(x0, y0, maxIterations, result)
	}
}

func calculateIterations(x0 *big.Float, y0 *big.Float, maxIterations int16, result *MandelbrotResult) int16 {
	nativeX0, _ := x0.Float64()
	nativeY0, _ := y0.Float64()
	return calculateIterationsWithDouble(nativeX0, nativeY0, maxIterations, result)
}

func calculateIterationsWithDouble(x0 float64, y0 float64, maxIterations int16, result *MandelbrotResult) int16 {
	x2 := 0.0
	y2 := 0.0
	iterations := int16(0)
	x := 0.0
	y := 0.0
	for {
		if iterations >= maxIterations {
			result.IncNumInfinite()
			iterations = 0
			break
		}
		if x2+y2 > 4.0 {
			result.IncNumEscaped()
			break
		}
		y = (x+x)*y + y0
		x = x2 - y2 + x0
		x2 = x * x
		y2 = y * y
		iterations++
	}
	return iterations
}
