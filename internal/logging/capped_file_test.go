package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCappedFileKeepsOnlyLastBytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.txt")
	writer := NewCappedFile(path, 5)

	if _, err := writer.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte(" world")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "world" {
		t.Fatalf("log = %q, want %q", data, "world")
	}
}
