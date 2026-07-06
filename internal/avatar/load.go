package avatar

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (g *Generator) load(ctx context.Context, artURL string) (image.Image, error) {
	parsed, err := url.Parse(strings.TrimSpace(artURL))
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "file" || parsed.Scheme == "" {
		return decodeFile(parsed)
	}
	if parsed.Scheme == "http" || parsed.Scheme == "https" {
		return g.decodeHTTP(ctx, artURL)
	}
	return nil, fmt.Errorf("unsupported art URL scheme: %s", parsed.Scheme)
}

func decodeFile(parsed *url.URL) (image.Image, error) {
	path, err := filePath(parsed)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func filePath(parsed *url.URL) (string, error) {
	path, err := url.PathUnescape(parsed.Path)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" {
		path, err = url.PathUnescape(parsed.String())
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func (g *Generator) decodeHTTP(ctx context.Context, artURL string) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artURL, nil)
	if err != nil {
		return nil, err
	}
	res, err := g.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("cover download failed: %s", res.Status)
	}
	img, _, err := image.Decode(res.Body)
	return img, err
}
