package daemon

import (
	"os"
	"path/filepath"
	"strconv"
)

type Paths struct {
	Dir             string
	PID             string
	Log             string
	ShutdownRestore string
}

func NewPaths(dir string) Paths {
	return Paths{
		Dir:             dir,
		PID:             filepath.Join(dir, "mpris-tg-status.pid"),
		Log:             filepath.Join(runtimeDir(), "mpris-tg-status.log"),
		ShutdownRestore: filepath.Join(dir, "shutdown_restore.txt"),
	}
}

func runtimeDir() string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return dir
	}
	return filepath.Join("/dev/shm", "mpris-tg-status-"+strconv.Itoa(os.Getuid()))
}
