package daemon

import (
	"os"
	"syscall"
)

func processAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}

func stopProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	waitUntilStopped(pid)
	return nil
}
