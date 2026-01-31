package mandelbrot

import (
	"math/big"
	"sync"
)

// const CalculationTypeBigDecimal int32 = 0
const CalculationTypeDouble int32 = 1

type MandelbrotResult struct {
	IterationCounts   [][]int16
	Width             int32
	Height            int32
	MaxIterations     int16
	NumEscaped        int32
	NumInfinite       int32
	CalculationTimeMs int64
	CoordMinX         big.Float
	CoordMaxX         big.Float
	CoordMinY         big.Float
	CoordMaxY         big.Float
	CalculationType   int32
	mutex             sync.Mutex
}

func (mr *MandelbrotResult) IncNumEscaped() {
	mr.NumEscaped++
}

func (mr *MandelbrotResult) IncNumInfinite() {
	mr.NumInfinite++
}

func NewMandelbrotResult(width int32, height int32, calculationType int32) *MandelbrotResult {
	result := MandelbrotResult{
		Width:           width,
		Height:          height,
		CalculationType: calculationType,
		IterationCounts: make([][]int16, width),
	}
	for i := range result.IterationCounts {
		result.IterationCounts[i] = make([]int16, height)
	}
	return &result
}
