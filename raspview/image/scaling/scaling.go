package scaling

import (
	"image"
	"image/color"
	"math"
)

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
