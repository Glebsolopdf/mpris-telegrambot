package status

import (
	"context"
	"errors"
	"time"

	statustext "mpris-tg-status/internal/status/text"
)

type Dependencies struct {
	Player TelegramPlayer
	Telegram
	Avatar Avatar
	Logger Logger
	Settings
	Clock    Clock
	Emoji    EmojiPicker
	Shutdown ShutdownPolicy
}

type Service struct {
	deps            Dependencies
	lastBio         string
	lastFirstName   string
	lastLastName    string
	lastAvatarKey   string
	avatarMode      avatarMode
	activeTrackID   string
	activeNameEmoji string
	activeBioEmoji  string
	nextText        time.Time
	nextAvatar      time.Time
	noPlayerMode    bool
}

func New(deps Dependencies) *Service {
	if deps.Emoji == nil {
		deps.Emoji = statustext.RandomEmojiPicker{}
	}
	return &Service{deps: deps}
}

func (s *Service) tick(ctx context.Context) error {
	track, err := s.deps.Player.Current(ctx)
	if err != nil {
		if errors.Is(err, ErrNoPlayer) {
			if !s.noPlayerMode {
				s.log("no active MPRIS player found, using slow scan interval")
			}
			s.noPlayerMode = true
			return s.applyIdle(ctx)
		}
		return err
	}
	if s.noPlayerMode {
		s.log("MPRIS player found, using normal poll interval")
	}
	s.noPlayerMode = false

	if !track.Playing {
		s.debug("no track is playing")
		return s.applyIdle(ctx)
	}

	if track.ID != s.activeTrackID {
		s.activeTrackID = track.ID
		baseEmoji := s.deps.Emoji.Pick()
		s.activeNameEmoji = statustext.ActiveEmoji(s.deps.Settings.ActiveNameEmojis, baseEmoji)
		s.activeBioEmoji = statustext.ActiveEmoji(s.deps.Settings.ActiveBioEmojis, baseEmoji)
		s.log("new track: %s -- %s", track.Artist, track.Title)
	}

	if s.textAllowed() {
		if err := s.applyActiveText(ctx, track); err != nil {
			return err
		}
	} else {
		s.debug("text cooldown active until %s", s.nextText.Format(time.RFC3339))
	}
	return s.applyAvatar(ctx, track)
}

func (s *Service) applyIdle(ctx context.Context) error {
	if s.textAllowed() {
		if err := s.applyName(ctx, s.deps.Settings.DefaultFirstName, s.deps.Settings.DefaultLastName); err != nil {
			return err
		}
		if err := s.applyBio(ctx, s.deps.Settings.IdleBio); err != nil {
			return err
		}
	}
	return s.applyDefaultAvatar(ctx)
}

func (s *Service) textAllowed() bool {
	return !s.deps.Clock.Now().Before(s.nextText)
}

func (s *Service) pauseTextOnRetry(err error) {
	if delay := retryDelay(err); delay > 0 {
		s.nextText = s.deps.Clock.Now().Add(delay)
		s.log("telegram text flood control: retry after %s", delay)
	}
}

func (s *Service) log(format string, v ...any) {
	if s.deps.Logger != nil {
		s.deps.Logger.Printf(format, v...)
	}
}

func (s *Service) debug(format string, v ...any) {
	if logger, ok := s.deps.Logger.(DebugLogger); ok {
		logger.Debugf(format, v...)
	}
}

type TelegramPlayer interface {
	Player
}

var ErrNoPlayer = errors.New("no active MPRIS player found")
