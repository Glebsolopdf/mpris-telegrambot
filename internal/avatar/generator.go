package avatar

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"net/http"
)

const (
	canvasSize  = 640
	coverSize   = 440
	jpegQuality = 92
)

type Generator struct {
	http *http.Client
}

func NewGenerator(httpClient *http.Client) *Generator {
	return &Generator{http: httpClient}
}

func (g *Generator) Build(ctx context.Context, artURL string, _ string) ([]byte, error) {
	source, err := g.load(ctx, artURL)
	if err != nil {
		return nil, err
	}

	square := centerSquare(source)
	background := resizeNearest(square, canvasSize, canvasSize)
	background = boxBlur(background, 18, 3)

	cover := resizeNearest(square, coverSize, coverSize)
	canvas := image.NewRGBA(image.Rect(0, 0, canvasSize, canvasSize))
	drawImage(canvas, background, 0, 0)
	drawShadow(canvas)
	drawImage(canvas, cover, (canvasSize-coverSize)/2, (canvasSize-coverSize)/2)

	var out bytes.Buffer
	if err := jpeg.Encode(&out, canvas, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (g *Generator) LoadDefault(ctx context.Context, path string) ([]byte, error) {
	source, err := g.load(ctx, path)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := jpeg.Encode(&out, source, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
