package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func currentExecutable() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(executable); err == nil {
		return resolved
	}
	return executable
}

func ensureMSAlias(enabled bool, executable string) (string, error) {
	if !enabled {
		return "", nil
	}
	if executable == "" {
		return "", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	aliasPath := filepath.Join(home, ".local", "bin", "ms")
	if aliasPointsTo(aliasPath, executable) {
		return "ms alias already configured", nil
	}
	if info, err := os.Lstat(aliasPath); err == nil {
		if !replaceableAlias(aliasPath, info) {
			return fmt.Sprintf("ms alias already exists at %s and was left unchanged", aliasPath), nil
		}
		if err := os.Remove(aliasPath); err != nil {
			return "", err
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(aliasPath), 0o700); err != nil {
		return "", err
	}
	if err := os.Symlink(executable, aliasPath); err != nil {
		return "", err
	}
	return fmt.Sprintf("ms alias configured at %s", aliasPath), nil
}

func aliasPointsTo(aliasPath string, executable string) bool {
	target, err := filepath.EvalSymlinks(aliasPath)
	if err != nil {
		return false
	}
	return filepath.Clean(target) == filepath.Clean(executable)
}

func replaceableAlias(aliasPath string, info os.FileInfo) bool {
	if info.Mode()&os.ModeSymlink == 0 {
		return false
	}
	target, err := os.Readlink(aliasPath)
	if err != nil {
		return false
	}
	return filepath.Base(target) == "mpris-tg-status"
}
