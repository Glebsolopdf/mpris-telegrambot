package daemon

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m Manager) DeleteData() (string, error) {
	stopMessage, err := m.DownWithRestore(true)
	if err != nil {
		return "", err
	}

	removed := 0
	for _, path := range m.deletePaths() {
		if removePath(path) {
			removed++
		}
	}
	for _, path := range m.aliasPaths() {
		if removeAlias(path, m.binaryPaths()) {
			removed++
		}
	}
	return fmt.Sprintf("%s\ndeleted %d local files", stopMessage, removed), nil
}

func (m Manager) deletePaths() []string {
	paths := []string{m.Paths.Log, autostartPath()}
	for _, dir := range m.dataDirs() {
		paths = append(paths,
			filepath.Join(dir, ".env"),
			filepath.Join(dir, "business_connection.json"),
			filepath.Join(dir, "avatar_cooldown.json"),
			filepath.Join(dir, "mpris-tg-status.pid"),
			filepath.Join(dir, "shutdown_restore.txt"),
		)
	}
	return paths
}

func (m Manager) dataDirs() []string {
	dirs := []string{m.Paths.Dir}
	if filepath.Base(m.Paths.Dir) == "bin" {
		dirs = append(dirs, filepath.Dir(m.Paths.Dir))
	} else {
		dirs = append(dirs, filepath.Join(m.Paths.Dir, "bin"))
	}
	return dirs
}

func (m Manager) aliasPaths() []string {
	paths := []string{
		filepath.Join(m.Paths.Dir, "ms"),
	}
	if home := os.Getenv("HOME"); home != "" {
		paths = append(paths, filepath.Join(home, ".local", "bin", "ms"))
	}
	if filepath.Base(m.Paths.Dir) == "bin" {
		paths = append(paths, filepath.Join(filepath.Dir(m.Paths.Dir), "ms"))
	}
	return paths
}

func (m Manager) binaryPaths() []string {
	paths := []string{
		m.Executable,
		filepath.Join(m.Paths.Dir, "mpris-tg-status"),
		filepath.Join(m.Paths.Dir, "bin", "mpris-tg-status"),
	}
	if filepath.Base(m.Paths.Dir) == "bin" {
		parent := filepath.Dir(m.Paths.Dir)
		paths = append(paths, filepath.Join(parent, "mpris-tg-status"))
	}
	return paths
}

func autostartPath() string {
	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "autostart", "mpris-tg-status.desktop")
}

func removePath(path string) bool {
	if path == "" {
		return false
	}
	if err := os.Remove(path); err == nil {
		return true
	}
	return false
}

func removeAlias(path string, binaries []string) bool {
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}
	for _, binary := range binaries {
		if samePath(target, binary) {
			return removePath(path)
		}
	}
	return false
}

func samePath(left string, right string) bool {
	return filepath.Clean(left) == filepath.Clean(right)
}
