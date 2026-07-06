package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
	Executable string
	Paths      Paths
}

func (m Manager) Up() (string, error) {
	return m.UpWithDelay(0)
}

func (m Manager) UpWithDelay(delay time.Duration) (string, error) {
	if pid, err := readPID(m.Paths.PID); err == nil && processRunning(pid) {
		return fmt.Sprintf("already running, pid=%d", pid), nil
	}
	_ = os.Remove(m.Paths.PID)
	_ = os.Remove(m.Paths.ShutdownRestore)

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		return "", err
	}
	defer devNull.Close()

	if err := os.MkdirAll(filepath.Dir(m.Paths.Log), 0o700); err != nil {
		return "", err
	}
	env := os.Environ()
	env = append(env, "MPRIS_TG_LOG_PATH="+m.Paths.Log)
	if delay > 0 {
		env = append(env, "MPRIS_TG_STARTUP_DELAY="+delay.String())
	}
	process, err := os.StartProcess(m.Executable, []string{m.Executable, "run"}, &os.ProcAttr{
		Dir:   m.Paths.Dir,
		Env:   env,
		Files: []*os.File{devNull, devNull, devNull},
		Sys:   processAttrs(),
	})
	if err != nil {
		return "", err
	}
	if err := writePID(m.Paths.PID, process.Pid); err != nil {
		_ = process.Signal(os.Interrupt)
		return "", err
	}
	_ = process.Release()
	return fmt.Sprintf("started, pid=%d, log=%s", process.Pid, m.Paths.Log), nil
}

func (m Manager) Down() (string, error) {
	return m.DownWithRestore(true)
}

func (m Manager) DownWithRestore(restore bool) (string, error) {
	pid, err := readPID(m.Paths.PID)
	if errors.Is(err, os.ErrNotExist) {
		return "not running", nil
	}
	if err != nil {
		_ = os.Remove(m.Paths.PID)
		return "removed broken pid file", nil
	}
	if !processRunning(pid) {
		_ = os.Remove(m.Paths.PID)
		return "not running, removed stale pid file", nil
	}
	if err := WriteShutdownRestore(m.Paths.ShutdownRestore, restore); err != nil {
		return "", err
	}
	if err := stopProcess(pid); err != nil {
		return "", err
	}
	if processRunning(pid) {
		return "", fmt.Errorf("process %d did not stop after timeout", pid)
	}
	_ = os.Remove(m.Paths.PID)
	return fmt.Sprintf("stopped, pid=%d", pid), nil
}

func (m Manager) Restart() (string, error) {
	return m.RestartWithOptions(true, 0)
}

func (m Manager) RestartWithOptions(restore bool, delay time.Duration) (string, error) {
	down, err := m.DownWithRestore(restore)
	if err != nil {
		return "", err
	}
	up, err := m.UpWithDelay(delay)
	if err != nil {
		return "", err
	}
	return down + "\n" + up, nil
}

func (m Manager) Status() string {
	pid, err := readPID(m.Paths.PID)
	if err != nil || !processRunning(pid) {
		return "down"
	}
	return fmt.Sprintf("up, pid=%d, log=%s", pid, m.Paths.Log)
}

func waitUntilStopped(pid int) {
	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) && processRunning(pid) {
		time.Sleep(100 * time.Millisecond)
	}
}
