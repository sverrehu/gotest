package mandelbrot

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
)

const version = 4

func Save(result *MandelbrotResult, fileName string) error {
	data, err := SaveToByteArray(result)
	if err != nil {
		return err
	}
	fmt.Printf("Saving %d bytes to %s\n", len(data), fileName)
	return os.WriteFile(fileName, data, 0644)
}

//goland:noinspection GoRedundantConversion,GoUnhandledErrorResult
func SaveToByteArray(result *MandelbrotResult) ([]byte, error) {
	endianness := binary.BigEndian
	var data bytes.Buffer
	data.WriteByte(byte(version))
	binary.Write(&data, endianness, int32(result.Width))
	binary.Write(&data, endianness, int32(result.Height))
	writeBigDecimal(&data, result.CoordMinX)
	writeBigDecimal(&data, result.CoordMaxX)
	writeBigDecimal(&data, result.CoordMinY)
	writeBigDecimal(&data, result.CoordMaxY)
	binary.Write(&data, endianness, int32(result.MaxIterations))
	binary.Write(&data, endianness, int32(result.NumEscaped))
	binary.Write(&data, endianness, int32(result.NumInfinite))
	binary.Write(&data, endianness, int64(result.CalculationTimeMs))
	binary.Write(&data, endianness, int32(-1 /* result.MathPrecision */))
	binary.Write(&data, endianness, int32(result.CalculationType))
	for x := range result.IterationCounts {
		for y := range result.IterationCounts[x] {
			binary.Write(&data, endianness, int16(result.IterationCounts[x][y]))
		}
	}
	return compress(data.Bytes())
}

//goland:noinspection GoUnhandledErrorResult
func writeBigDecimal(data *bytes.Buffer, n big.Float) {
	endianness := binary.BigEndian
	s := n.Text('g', -1)
	binary.Write(data, endianness, int32(len(s)))
	for _, c := range s {
		binary.Write(data, endianness, int8(c))
	}
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	defer func(gz *gzip.Writer) {
		_ = gz.Close()
	}(gz)
	if err != nil {
		return nil, fmt.Errorf("error getting gzip writer: %w", err)
	}
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	return buf.Bytes(), nil
}
