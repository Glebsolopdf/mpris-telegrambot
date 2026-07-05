package logging

import (
	"os"
	"sync"
)

type CappedFile struct {
	path string
	max  int
	mu   sync.Mutex
}

func NewCappedFile(path string, max int) *CappedFile {
	return &CappedFile{path: path, max: max}
}

func (w *CappedFile) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := os.ReadFile(w.path)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	next := append(data, p...)
	if len(next) > w.max {
		next = next[len(next)-w.max:]
	}
	if err := os.WriteFile(w.path, next, 0o600); err != nil {
		return 0, err
	}
	return len(p), nil
}
