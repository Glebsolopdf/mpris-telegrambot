package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromDirCreatesEnv(t *testing.T) {
	dir := t.TempDir()

	_, err := LoadFromDir(dir)
	if !errors.Is(err, ErrEnvCreated) {
		t.Fatalf("err = %v, want ErrEnvCreated", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".env")); err != nil {
		t.Fatal(err)
	}
}

func TestLoadUsesOptionalDurations(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := strings.Join([]string{
		"TELEGRAM_BOT_TOKEN=token",
		"TELEGRAM_TARGET_USER_ID=42",
		"POLL_INTERVAL=5s",
		"NO_PLAYER_POLL_INTERVAL=3m",
		"TELEGRAM_MIN_UPDATE_INTERVAL=2m",
		"HTTP_TIMEOUT=3s",
		"LOG_LEVEL=debug",
		"GENERATED_AVATAR_ENABLED=false",
		"MEMORY_LIMIT_MB=64",
		"GC_PERCENT=50",
		"STARTUP_DELAY=5m",
		"PROMPT_TIMEOUT=10s",
		"BUSINESS_CONNECTION_WAIT_INTERVAL=30s",
		"USE_MS_ALIAS=true",
		"ACTIVE_BIO_TEMPLATE={emoji_bio} Слушаю: {artist}",
		"ACTIVE_FIRST_NAME_TEMPLATE={default_first_name} {emoji_name}",
		"ACTIVE_NAME_EMOJIS=🔥",
		"ACTIVE_BIO_EMOJIS=false",
	}, "\n") + "\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TargetUserID != 42 {
		t.Fatalf("target user id = %d", cfg.TargetUserID)
	}
	if cfg.StatusSettings.PollInterval.String() != "5s" {
		t.Fatalf("poll interval = %s", cfg.StatusSettings.PollInterval)
	}
	if cfg.StatusSettings.NoPlayerPollInterval.String() != "3m0s" {
		t.Fatalf("no player poll interval = %s", cfg.StatusSettings.NoPlayerPollInterval)
	}
	if cfg.StatusSettings.TelegramMinInterval.String() != "2m0s" {
		t.Fatalf("telegram interval = %s", cfg.StatusSettings.TelegramMinInterval)
	}
	if cfg.HTTPTimeout.String() != "3s" {
		t.Fatalf("http timeout = %s", cfg.HTTPTimeout)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("log level = %s", cfg.LogLevel)
	}
	if !cfg.StatusSettings.DisableGeneratedAvatar {
		t.Fatal("generated avatar updates should be disabled")
	}
	if cfg.MemoryLimitMB != 64 || cfg.GCPercent != 50 {
		t.Fatalf("memory settings = %d MB GOGC %d", cfg.MemoryLimitMB, cfg.GCPercent)
	}
	if cfg.StartupDelay.String() != "5m0s" || cfg.PromptTimeout.String() != "10s" {
		t.Fatalf("prompt settings = %s/%s", cfg.StartupDelay, cfg.PromptTimeout)
	}
	if cfg.BusinessWaitInterval.String() != "30s" {
		t.Fatalf("business wait interval = %s", cfg.BusinessWaitInterval)
	}
	if !cfg.UseMSAlias {
		t.Fatal("ms alias should be enabled")
	}
	if cfg.StatusSettings.ActiveBioTemplate != "{emoji_bio} Слушаю: {artist}" {
		t.Fatalf("active bio template = %q", cfg.StatusSettings.ActiveBioTemplate)
	}
	if cfg.StatusSettings.ActiveFirstNameTemplate != "{default_first_name} {emoji_name}" {
		t.Fatalf("active first name template = %q", cfg.StatusSettings.ActiveFirstNameTemplate)
	}
	if cfg.StatusSettings.ActiveNameEmojis != "🔥" || cfg.StatusSettings.ActiveBioEmojis != "false" {
		t.Fatalf("active emojis = %q/%q", cfg.StatusSettings.ActiveNameEmojis, cfg.StatusSettings.ActiveBioEmojis)
	}
}
