package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func env(values map[string]string, name string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return strings.TrimSpace(values[name])
}

func int64Env(values map[string]string, name string) int64 {
	parsed, _ := strconv.ParseInt(env(values, name), 10, 64)
	return parsed
}

func intEnv(values map[string]string, name string, fallback int) int {
	value := env(values, name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return parsed
}

func durationEnv(values map[string]string, name string, fallback time.Duration) time.Duration {
	value := env(values, name)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return -1
	}
	return parsed
}

func boolEnv(values map[string]string, name string, fallback bool) bool {
	value := strings.ToLower(env(values, name))
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "yes", "on", "enabled":
		return true
	case "0", "false", "no", "off", "disabled":
		return false
	default:
		return fallback
	}
}

func defaultBio(values map[string]string) string {
	if value := env(values, "DEFAULT_BIO"); value != "" {
		return value
	}
	return env(values, "IDLE_BIO")
}

func logLevel(values map[string]string) string {
	value := env(values, "LOG_LEVEL")
	if value == "" {
		return "info"
	}
	return strings.ToLower(value)
}
