package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"mpris-tg-status/internal/avatar"
	"mpris-tg-status/internal/bootstrap"
	"mpris-tg-status/internal/business"
	"mpris-tg-status/internal/config"
	"mpris-tg-status/internal/daemon"
	"mpris-tg-status/internal/logging"
	"mpris-tg-status/internal/mpris"
	"mpris-tg-status/internal/status"
	"mpris-tg-status/internal/telegram"
)

func runService() error {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrEnvCreated) {
			return err
		}
		return errors.New("config: " + err.Error())
	}
	debug.SetGCPercent(cfg.GCPercent)
	debug.SetMemoryLimit(int64(cfg.MemoryLimitMB) * 1024 * 1024)

	logger, err := newLogger(cfg.LogLevel)
	if err != nil {
		return err
	}
	logger.Printf("config loaded from %s", cfg.EnvPath)
	logger.Printf("memory tuned: limit=%d MB gc_percent=%d", cfg.MemoryLimitMB, cfg.GCPercent)
	if message, err := ensureMSAlias(cfg.UseMSAlias, currentExecutable()); err != nil {
		logger.Printf("ms alias setup failed: %v", err)
	} else if message != "" {
		logger.Printf(message)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if delay := requestedStartupDelay(); delay > 0 {
		logger.Printf("startup delay before external requests: %s", delay)
		if err := waitBeforeStart(ctx, delay); err != nil {
			return err
		}
	}

	return startService(ctx, cfg, logger)
}

func startService(ctx context.Context, cfg config.Config, logger logging.Logger) error {
	httpClient := &http.Client{Timeout: cfg.HTTPTimeout}
	for {
		err := runConnectedService(ctx, cfg, httpClient, logger)
		if !status.IsBusinessConnectionInvalid(err) {
			return err
		}
		logger.Printf("business connection id became invalid, waiting for a fresh id")
		if err := business.NewStore(cfg.BusinessStatePath).Delete(); err != nil {
			return err
		}
	}
}

func runConnectedService(ctx context.Context, cfg config.Config, httpClient *http.Client, logger logging.Logger) error {
	businessState, err := loadBusinessState(ctx, cfg, httpClient, logger)
	if err != nil {
		return errors.New("telegram business setup: " + err.Error())
	}

	cfg.StatusSettings.DefaultFirstName = firstNonEmpty(cfg.StatusSettings.DefaultFirstName, businessState.FirstName)
	cfg.StatusSettings.DefaultLastName = firstNonEmpty(cfg.StatusSettings.DefaultLastName, businessState.LastName)
	logger.Printf("business user id=%d default_name=%q %q", cfg.TargetUserID, cfg.StatusSettings.DefaultFirstName, cfg.StatusSettings.DefaultLastName)

	mprisClient, err := mpris.NewClient(cfg.PreferredPlayer)
	if err != nil {
		return errors.New("mpris: " + err.Error())
	}
	defer mprisClient.Close()
	logger.Printf("mpris connected preferred_player=%q", cfg.PreferredPlayer)

	service := status.New(status.Dependencies{
		Player:   mprisClient,
		Telegram: telegram.NewClient(cfg.TelegramBotToken, businessState.ConnectionID, httpClient),
		Avatar:   avatar.NewGenerator(httpClient),
		Logger:   logger,
		Settings: cfg.StatusSettings,
		Clock:    status.RealClock{},
		Shutdown: daemon.ShutdownPolicy{Path: cfg.ShutdownRestorePath},
	})
	return service.Run(ctx)
}

func loadBusinessState(ctx context.Context, cfg config.Config, httpClient *http.Client, logger logging.Logger) (business.State, error) {
	setupTelegram := telegram.NewClient(cfg.TelegramBotToken, "", httpClient)
	stateStore := business.NewStore(cfg.BusinessStatePath)
	return bootstrap.BusinessConnectionID(ctx, stateStore, setupTelegram, cfg.TargetUserID, cfg.BusinessWaitInterval, logger)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
