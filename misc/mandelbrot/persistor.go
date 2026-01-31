package mandelbrot

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"os"
)

const version = 3

func Save(result *MandelbrotResult, fileName string) error {
	data, err := SaveToByteArray(result)
	if err == nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func SaveToByteArray(result *MandelbrotResult) ([]byte, error) {
	endianness := binary.BigEndian
	var data bytes.Buffer
	data.WriteByte(byte(version))
	binary.Write(&data, endianness, int32(result.Width))
	binary.Write(&data, endianness, int32(result.Height))
	writeBigDecimal(data, result.CoordMinX)
	writeBigDecimal(data, result.CoordMaxX)
	writeBigDecimal(data, result.CoordMinY)
	writeBigDecimal(data, result.CoordMaxY)
	binary.Write(&data, endianness, int32(result.MaxIterations))
	binary.Write(&data, endianness, int32(result.NumEscaped))
	binary.Write(&data, endianness, int32(result.NumInfinite))
	binary.Write(&data, endianness, int64(result.CalculationTimeMs))
	binary.Write(&data, endianness, int32(11 /* result.MathPrecision */))
	binary.Write(&data, endianness, int32(result.CalculationType))
	for x, _ := range result.IterationCounts {
		for y, _ := range result.IterationCounts[x] {
			binary.Write(&data, endianness, int16(result.IterationCounts[x][y]))
		}
	}

	// TODO
	return data.Bytes(), nil
}

func writeBigDecimal(data bytes.Buffer, n big.Float) {

}
