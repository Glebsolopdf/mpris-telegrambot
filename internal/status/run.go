package status

import (
	"context"
	"time"
)

const shutdownRestoreTimeout = 30 * time.Second

func (s *Service) Run(ctx context.Context) error {
	s.log(
		"service started: poll=%s no_player_poll=%s telegram_min_interval=%s",
		s.deps.Settings.PollInterval,
		s.noPlayerPollInterval(),
		s.deps.Settings.TelegramMinInterval,
	)
	if err := s.tick(ctx); err != nil {
		if IsBusinessConnectionInvalid(err) {
			return err
		}
		s.log("initial tick: %v", err)
	}

	timer := time.NewTimer(s.nextPollInterval())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			s.restoreOnShutdown()
			s.log("service stopped")
			return nil
		case <-timer.C:
			if err := s.tick(ctx); err != nil {
				if IsBusinessConnectionInvalid(err) {
					return err
				}
				s.log("tick: %v", err)
			}
			timer.Reset(s.nextPollInterval())
		}
	}
}

func (s *Service) nextPollInterval() time.Duration {
	if s.noPlayerMode {
		return s.noPlayerPollInterval()
	}
	return s.deps.Settings.PollInterval
}

func (s *Service) noPlayerPollInterval() time.Duration {
	if s.deps.Settings.NoPlayerPollInterval > 0 {
		return s.deps.Settings.NoPlayerPollInterval
	}
	return s.deps.Settings.PollInterval
}

func (s *Service) restoreOnShutdown() {
	if s.deps.Shutdown != nil && !s.deps.Shutdown.RestoreProfile() {
		s.log("default profile restore skipped before shutdown")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownRestoreTimeout)
	defer cancel()

	s.log("restoring default profile before shutdown")
	s.nextText = time.Time{}
	s.nextAvatar = time.Time{}
	if err := s.applyIdle(ctx); err != nil {
		s.log("shutdown restore failed: %v", err)
	}
}
