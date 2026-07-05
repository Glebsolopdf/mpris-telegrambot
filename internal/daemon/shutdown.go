package daemon

import (
	"os"
	"strings"
)

type ShutdownPolicy struct {
	Path string
}

func WriteShutdownRestore(path string, restore bool) error {
	value := "restore\n"
	if !restore {
		value = "skip\n"
	}
	return os.WriteFile(path, []byte(value), 0o600)
}

func (p ShutdownPolicy) RestoreProfile() bool {
	if p.Path == "" {
		return true
	}

	data, err := os.ReadFile(p.Path)
	_ = os.Remove(p.Path)
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(data)) != "skip"
}
