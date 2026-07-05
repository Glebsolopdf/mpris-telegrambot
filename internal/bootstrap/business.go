package bootstrap

import (
	"context"
	"errors"
	"time"

	"mpris-tg-status/internal/business"
	"mpris-tg-status/internal/telegram"
)

type Logger interface {
	Printf(format string, v ...any)
}

func BusinessConnectionID(
	ctx context.Context,
	store business.Store,
	client *telegram.Client,
	userID int64,
	retryInterval time.Duration,
	logger Logger,
) (business.State, error) {
	for {
		state, ok, err := store.LoadForUser(userID)
		if err != nil {
			return business.State{}, err
		}
		if ok {
			if logger != nil {
				logger.Printf("loaded Telegram business connection id for user %d", userID)
			}
			return state, nil
		}

		if logger != nil {
			logger.Printf("searching Telegram business connection id for user %d", userID)
		}
		connection, err := client.FindBusinessConnection(ctx, userID)
		if err == nil {
			return saveConnection(store, connection, logger)
		}
		if !isPendingConnection(err) {
			return business.State{}, err
		}
		if logger != nil {
			logger.Printf("business connection id is unavailable, waiting %s", retryInterval)
		}
		if err := wait(ctx, retryInterval); err != nil {
			return business.State{}, err
		}
	}
}

func isPendingConnection(err error) bool {
	return errors.Is(err, telegram.ErrBusinessConnectionNotFound) ||
		errors.Is(err, telegram.ErrBusinessConnectionDisabled)
}

func saveConnection(store business.Store, connection telegram.BusinessConnection, logger Logger) (business.State, error) {
	state := business.State{
		ConnectionID: connection.ID,
		UserID:       connection.UserID,
		FirstName:    connection.FirstName,
		LastName:     connection.LastName,
	}
	if err := store.Save(state); err != nil {
		return business.State{}, err
	}
	if logger != nil {
		logger.Printf("saved Telegram business connection id for user %d", connection.UserID)
	}
	return state, nil
}

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
