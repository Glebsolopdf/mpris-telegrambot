package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"mpris-tg-status/internal/logging"
)

const logMaxBytes = 1024 * 1024

func newLogger(level string) (logging.Logger, error) {
	path := os.Getenv("MPRIS_TG_LOG_PATH")
	if path == "" {
		return logging.New(level), nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return logging.Logger{}, err
	}
	return logging.New(level, logging.NewCappedFile(path, logMaxBytes)), nil
}

func requestedStartupDelay() time.Duration {
	value := os.Getenv("MPRIS_TG_STARTUP_DELAY")
	delay, err := time.ParseDuration(value)
	if err != nil || delay < 0 {
		return 0
	}
	return delay
}

func waitBeforeStart(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
