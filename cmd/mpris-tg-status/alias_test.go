package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureMSAliasReplacesOldBinarySymlink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	aliasPath := filepath.Join(home, ".local", "bin", "ms")
	if err := os.MkdirAll(filepath.Dir(aliasPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/old/path/mpris-tg-status", aliasPath); err != nil {
		t.Fatal(err)
	}

	executable := filepath.Join(t.TempDir(), "mpris-tg-status")
	message, err := ensureMSAlias(true, executable)
	if err != nil {
		t.Fatal(err)
	}
	if message == "" {
		t.Fatal("expected setup message")
	}
	target, err := os.Readlink(aliasPath)
	if err != nil {
		t.Fatal(err)
	}
	if target != executable {
		t.Fatalf("alias target = %q, want %q", target, executable)
	}
}
