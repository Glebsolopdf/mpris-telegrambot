package avatar

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestBuildCreatesSquareJPEG(t *testing.T) {
	path := writeTestCover(t, "cover-*.jpg")

	jpegBytes, err := NewGenerator(&http.Client{}).Build(context.Background(), "file://"+path, "Track")
	if err != nil {
		t.Fatal(err)
	}

	result, err := jpeg.Decode(bytes.NewReader(jpegBytes))
	if err != nil {
		t.Fatal(err)
	}
	if result.Bounds().Dx() != canvasSize || result.Bounds().Dy() != canvasSize {
		t.Fatalf("size = %v", result.Bounds())
	}
}

func TestBuildReadsEscapedFileURL(t *testing.T) {
	path := writeTestCover(t, "обложка *.jpg")
	escapedURL := (&url.URL{Scheme: "file", Path: path}).String()

	if _, err := NewGenerator(&http.Client{}).LoadDefault(context.Background(), escapedURL); err != nil {
		t.Fatal(err)
	}
}

func TestBuildDoesNotRenderTrackTitle(t *testing.T) {
	path := writeTestCover(t, "cover-*.jpg")
	generator := NewGenerator(&http.Client{})

	first, err := generator.Build(context.Background(), "file://"+path, "First Track")
	if err != nil {
		t.Fatal(err)
	}
	second, err := generator.Build(context.Background(), "file://"+path, "Another Track")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("avatar changed when only track title changed")
	}
}

func writeTestCover(t *testing.T, pattern string) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), pattern)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, 48, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 48; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 6), B: 120, A: 255})
		}
	}
	if err := jpeg.Encode(file, img, nil); err != nil {
		t.Fatal(err)
	}
	return file.Name()
}
