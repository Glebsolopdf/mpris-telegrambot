package status

import (
	"context"
	"time"
)

type Settings struct {
	PollInterval            time.Duration
	NoPlayerPollInterval    time.Duration
	TelegramMinInterval     time.Duration
	AvatarMinInterval       time.Duration
	DisableGeneratedAvatar  bool
	IdleBio                 string
	ActiveBioTemplate       string
	ActiveFirstNameTemplate string
	ActiveNameEmojis        string
	ActiveBioEmojis         string
	DefaultAvatarPath       string
	DefaultFirstName        string
	DefaultLastName         string
}

type Track struct {
	ID      string
	Artist  string
	Title   string
	ArtURL  string
	Playing bool
}

type Player interface {
	Current(ctx context.Context) (Track, error)
}

type Telegram interface {
	SetBio(ctx context.Context, bio string) error
	SetName(ctx context.Context, firstName string, lastName string) error
	SetAvatar(ctx context.Context, jpeg []byte) error
	RemoveAvatar(ctx context.Context) error
	RemovePublicAvatar(ctx context.Context) error
}

type Avatar interface {
	Build(ctx context.Context, artURL string, title string) ([]byte, error)
	LoadDefault(ctx context.Context, path string) ([]byte, error)
}

type ShutdownPolicy interface {
	RestoreProfile() bool
}

type Logger interface {
	Printf(format string, v ...any)
}

type DebugLogger interface {
	Debugf(format string, v ...any)
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}
