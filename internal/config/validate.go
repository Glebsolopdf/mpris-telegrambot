package config

import "errors"

func validate(cfg Config) error {
	if cfg.TelegramBotToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.TargetUserID <= 0 {
		return errors.New("TELEGRAM_TARGET_USER_ID must be a positive integer")
	}
	if cfg.StatusSettings.PollInterval <= 0 {
		return errors.New("POLL_INTERVAL must be positive")
	}
	if cfg.StatusSettings.NoPlayerPollInterval <= 0 {
		return errors.New("NO_PLAYER_POLL_INTERVAL must be positive")
	}
	if cfg.StatusSettings.TelegramMinInterval <= 0 {
		return errors.New("TELEGRAM_MIN_UPDATE_INTERVAL must be positive")
	}
	if cfg.StatusSettings.AvatarMinInterval <= 0 {
		return errors.New("AVATAR_MIN_UPDATE_INTERVAL must be positive")
	}
	if cfg.HTTPTimeout <= 0 {
		return errors.New("HTTP_TIMEOUT must be positive")
	}
	if cfg.MemoryLimitMB <= 0 {
		return errors.New("MEMORY_LIMIT_MB must be positive")
	}
	if cfg.GCPercent <= 0 {
		return errors.New("GC_PERCENT must be positive")
	}
	if cfg.StartupDelay < 0 {
		return errors.New("STARTUP_DELAY cannot be negative")
	}
	if cfg.PromptTimeout <= 0 {
		return errors.New("PROMPT_TIMEOUT must be positive")
	}
	if cfg.BusinessWaitInterval <= 0 {
		return errors.New("BUSINESS_CONNECTION_WAIT_INTERVAL must be positive")
	}
	return nil
}
