package status

import (
	"context"
	"testing"
	"time"
)

func TestRunRestoresDefaultProfileOnShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tg := &fakeTelegram{}
	service := New(Dependencies{
		Player: fakePlayer{track: Track{
			ID: "1", Artist: "A", Title: "T", ArtURL: "file:///cover.jpg", Playing: true,
		}},
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{
			PollInterval:        time.Hour,
			DefaultAvatarPath:   "default.jpg",
			DefaultFirstName:    "Name",
			DefaultLastName:     "Last",
			IdleBio:             "Default bio",
			TelegramMinInterval: time.Second,
		},
		Clock: fakeClock{now: time.Unix(10, 0)},
		Emoji: fixedEmoji{},
	})

	if err := service.Run(ctx); err != nil {
		t.Fatal(err)
	}
	if len(tg.bios) != 2 || tg.bios[1] != "Default bio" {
		t.Fatalf("bios = %v", tg.bios)
	}
	if len(tg.names) != 2 || tg.names[1] != "Name Last" {
		t.Fatalf("names = %v", tg.names)
	}
	if tg.avatars != 2 || tg.removals != 2 {
		t.Fatalf("avatars = %d removals = %d, want 2/2", tg.avatars, tg.removals)
	}
}

func TestRunSkipsDefaultProfileRestoreOnShutdown(t *testing.T) {
	tg := &fakeTelegram{}
	service := New(Dependencies{
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{
			DefaultAvatarPath: "default.jpg",
			DefaultFirstName:  "Name",
			IdleBio:           "Default bio",
		},
		Clock:    fakeClock{now: time.Unix(10, 0)},
		Shutdown: fakeShutdown{restore: false},
	})

	service.restoreOnShutdown()

	if len(tg.bios) != 0 || len(tg.names) != 0 || tg.avatars != 0 {
		t.Fatalf("updates = %d bios, %d names, %d avatars", len(tg.bios), len(tg.names), tg.avatars)
	}
}

func TestShutdownRestoreIgnoresAvatarCooldown(t *testing.T) {
	tg := &fakeTelegram{}
	service := New(Dependencies{
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{
			DefaultAvatarPath: "default.jpg",
			DefaultFirstName:  "Name",
			IdleBio:           "Default bio",
		},
		Clock: fakeClock{now: time.Unix(10, 0)},
	})
	service.nextAvatar = service.deps.Clock.Now().Add(time.Hour)

	service.restoreOnShutdown()

	if tg.avatars != 1 || tg.removals != 1 {
		t.Fatalf("avatars = %d removals = %d, want 1/1", tg.avatars, tg.removals)
	}
}

func TestNextPollIntervalSlowsDownWithoutPlayer(t *testing.T) {
	service := New(Dependencies{
		Settings: Settings{
			PollInterval:         20 * time.Second,
			NoPlayerPollInterval: 3 * time.Minute,
		},
	})

	if got := service.nextPollInterval(); got != 20*time.Second {
		t.Fatalf("normal poll interval = %s", got)
	}

	service.noPlayerMode = true
	if got := service.nextPollInterval(); got != 3*time.Minute {
		t.Fatalf("no player poll interval = %s", got)
	}
}

type fakeShutdown struct {
	restore bool
}

func (f fakeShutdown) RestoreProfile() bool {
	return f.restore
}
