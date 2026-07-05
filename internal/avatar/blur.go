package avatar

import (
	"image"
	"image/color"
)

func boxBlur(src *image.RGBA, radius int, passes int) *image.RGBA {
	out := src
	for i := 0; i < passes; i++ {
		out = blurVertical(blurHorizontal(out, radius), radius)
	}
	return out
}

func blurHorizontal(src *image.RGBA, radius int) *image.RGBA {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x, y, averageLine(src, x, y, radius, true))
		}
	}
	return out
}

func blurVertical(src *image.RGBA, radius int) *image.RGBA {
	bounds := src.Bounds()
	out := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x, y, averageLine(src, x, y, radius, false))
		}
	}
	return out
}

func averageLine(src *image.RGBA, cx int, cy int, radius int, horizontal bool) color.RGBA {
	var r, g, b, a, count uint32
	bounds := src.Bounds()
	for offset := -radius; offset <= radius; offset++ {
		x, y := cx, cy
		if horizontal {
			x += offset
		} else {
			y += offset
		}
		if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
			cr, cg, cb, ca := src.At(x, y).RGBA()
			r, g, b, a = r+cr, g+cg, b+cb, a+ca
			count++
		}
	}
	return color.RGBA{
		R: uint8((r / count) >> 8),
		G: uint8((g / count) >> 8),
		B: uint8((b / count) >> 8),
		A: uint8((a / count) >> 8),
	}
}
