package avatar

import (
	"image"
	"image/color"
)

func centerSquare(src image.Image) image.Image {
	bounds := src.Bounds()
	size := bounds.Dx()
	if bounds.Dy() < size {
		size = bounds.Dy()
	}

	x0 := bounds.Min.X + (bounds.Dx()-size)/2
	y0 := bounds.Min.Y + (bounds.Dy()-size)/2
	out := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			out.Set(x, y, src.At(x0+x, y0+y))
		}
	}
	return out
}

func resizeNearest(src image.Image, width int, height int) *image.RGBA {
	bounds := src.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		sy := bounds.Min.Y + y*bounds.Dy()/height
		for x := 0; x < width; x++ {
			sx := bounds.Min.X + x*bounds.Dx()/width
			out.Set(x, y, src.At(sx, sy))
		}
	}
	return out
}

func drawImage(dst *image.RGBA, src image.Image, dx int, dy int) {
	bounds := src.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			dst.Set(dx+x, dy+y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
		}
	}
}

func drawShadow(dst *image.RGBA) {
	start := (canvasSize - coverSize) / 2
	end := start + coverSize
	shadow := color.RGBA{A: 65}
	for y := start - 10; y < end+14; y++ {
		for x := start - 10; x < end+14; x++ {
			if x >= 0 && y >= 0 && x < canvasSize && y < canvasSize {
				dst.Set(x, y, blend(dst.At(x, y), shadow))
			}
		}
	}
}

func blend(base color.Color, overlay color.RGBA) color.RGBA {
	r, g, b, _ := base.RGBA()
	alpha := uint32(overlay.A)
	return color.RGBA{
		R: uint8((r*(255-alpha) + uint32(overlay.R)*257*alpha) / 255 / 257),
		G: uint8((g*(255-alpha) + uint32(overlay.G)*257*alpha) / 255 / 257),
		B: uint8((b*(255-alpha) + uint32(overlay.B)*257*alpha) / 255 / 257),
		A: 255,
	}
}
