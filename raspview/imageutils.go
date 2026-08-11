// Package imageutils contains helpers for encoding and transforming images.
package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"math"
)

// JPEGParams is the Go equivalent of the JPEGImageWriteParam used by the
// Java implementation. Quality is in the range [1, 100], as required by
// image/jpeg.Encode.
type JPEGParams struct {
	Quality int
}

var (
	JPEGParamsFast = JPEGParams{Quality: 69}
	JPEGParamsGood = JPEGParams{Quality: 80}
)

func JPEGEncodeFast(img image.Image) ([]byte, error) { return JPEGEncode(img, JPEGParamsFast) }

func JPEGEncodeGood(img image.Image) ([]byte, error) { return JPEGEncode(img, JPEGParamsGood) }

// JPEGEncode encodes img as JPEG. JPEG has no alpha channel, so transparent
// pixels are composited against black, matching the Java RGB conversion.
func JPEGEncode(img image.Image, params JPEGParams) ([]byte, error) {
	var out bytes.Buffer
	quality := params.Quality
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}
	if err := jpeg.Encode(&out, removeAlpha(img), &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// RemoveAlpha returns an opaque copy when img has an alpha channel. The
// flushOriginal argument is retained for source compatibility; Go images do
// not hold external resources that need flushing.
func RemoveAlpha(img image.Image, flushOriginal bool) image.Image {
	if !HasAlphaChannel(img) {
		return img
	}
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, b, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			dst.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}
	return dst
}

func removeAlpha(img image.Image) image.Image { return RemoveAlpha(img, false) }

// HasAlphaChannel reports whether img's concrete color model carries alpha.
func HasAlphaChannel(img image.Image) bool {
	switch img.(type) {
	case *image.RGBA, *image.RGBA64, *image.NRGBA, *image.NRGBA64,
		*image.Alpha, *image.Alpha16:
		return true
	case *image.Paletted:
		for _, c := range img.(*image.Paletted).Palette {
			_, _, _, a := c.RGBA()
			if a != 0xffff {
				return true
			}
		}
	}
	return false
}

// Decode decodes PNG, JPEG, or GIF data. It returns an error for malformed or
// unsupported data, where the Java method returned null or threw an unchecked
// I/O exception depending on the decoder.
func Decode(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// MakeRGBOrARGB is kept for parity with Java. image.Image already abstracts
// over concrete pixel formats, so no conversion is needed.
func MakeRGBOrARGB(img image.Image, flushOriginal bool) image.Image { return img }

func Copy(img image.Image) image.Image {
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			dst.Set(x, y, img.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return dst
}

func IsLandscape(img image.Image) bool { return img.Bounds().Dx() > img.Bounds().Dy() }

func LimitHeight(img image.Image, height int) image.Image {
	if img.Bounds().Dy() <= height {
		return img
	}
	return FitHeight(img, height)
}

func FitHeight(img image.Image, height int) image.Image {
	if img.Bounds().Dy() == height {
		return img
	}
	width := int(float64(img.Bounds().Dx()) * float64(height) / float64(img.Bounds().Dy()))
	return Resize(img, width, height)
}

func LimitWidth(img image.Image, width int) image.Image {
	if img.Bounds().Dx() <= width {
		return img
	}
	return FitWidth(img, width)
}

func FitWidth(img image.Image, width int) image.Image {
	if img.Bounds().Dx() == width {
		return img
	}
	height := int(float64(img.Bounds().Dy()) * float64(width) / float64(img.Bounds().Dx()))
	return Resize(img, width, height)
}

func Limit(img image.Image, maxWidth, maxHeight int) image.Image {
	widthRatio := float64(maxWidth) / float64(img.Bounds().Dx())
	heightRatio := float64(maxHeight) / float64(img.Bounds().Dy())
	if widthRatio < heightRatio {
		return LimitWidth(img, maxWidth)
	}
	return LimitHeight(img, maxHeight)
}

func Fit(img image.Image, width, height int) image.Image {
	widthRatio := float64(width) / float64(img.Bounds().Dx())
	heightRatio := float64(height) / float64(img.Bounds().Dy())
	if widthRatio < heightRatio {
		return FitWidth(img, width)
	}
	return FitHeight(img, height)
}

// Resize scales img using bilinear interpolation.
func Resize(img image.Image, width, height int) image.Image {
	if width < 1 || height < 1 {
		return image.NewRGBA(image.Rect(0, 0, max(width, 0), max(height, 0)))
	}
	src := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		sy := (float64(y)+0.5)*float64(src.Dy())/float64(height) - 0.5
		y0 := int(math.Floor(sy))
		fy := sy - float64(y0)
		for x := 0; x < width; x++ {
			sx := (float64(x)+0.5)*float64(src.Dx())/float64(width) - 0.5
			x0 := int(math.Floor(sx))
			fx := sx - float64(x0)
			c00 := rgbaAt(img, src, x0, y0)
			c10 := rgbaAt(img, src, x0+1, y0)
			c01 := rgbaAt(img, src, x0, y0+1)
			c11 := rgbaAt(img, src, x0+1, y0+1)
			dst.SetRGBA(x, y, bilinear(c00, c10, c01, c11, fx, fy))
		}
	}
	return dst
}

func CropToSquare(img image.Image) image.Image {
	b := img.Bounds()
	size := b.Dx()
	if b.Dy() < size {
		size = b.Dy()
	}
	x0, y0 := b.Min.X+(b.Dx()-size)/2, b.Min.Y+(b.Dy()-size)/2
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dst.Set(x, y, img.At(x0+x, y0+y))
		}
	}
	return dst
}

func GetFastParamsWithQuality(quality float32) JPEGParams {
	return JPEGParams{Quality: int(quality * 100)}
}
func GetGoodParamsWithQuality(quality float32) JPEGParams {
	return JPEGParams{Quality: int(quality * 100)}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rgbaAt(img image.Image, b image.Rectangle, x, y int) color.RGBA {
	if x < 0 {
		x = 0
	}
	if x >= b.Dx() {
		x = b.Dx() - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= b.Dy() {
		y = b.Dy() - 1
	}
	r, g, bl, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}
}

func bilinear(a, b, c, d color.RGBA, fx, fy float64) color.RGBA {
	lerp := func(v0, v1, t float64) float64 { return v0 + (v1-v0)*t }
	channel := func(v0, v1, v2, v3 uint8) uint8 {
		return uint8(math.Round(lerp(lerp(float64(v0), float64(v1), fx), lerp(float64(v2), float64(v3), fx), fy)))
	}
	return color.RGBA{channel(a.R, b.R, c.R, d.R), channel(a.G, b.G, c.G, d.G), channel(a.B, b.B, c.B, d.B), channel(a.A, b.A, c.A, d.A)}
}
