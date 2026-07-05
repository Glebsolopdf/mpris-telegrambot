package status

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTickAppliesTrackBioAndAvatarOnce(t *testing.T) {
	tg := &fakeTelegram{removeLimit: 1}
	service := New(Dependencies{
		Player: fakePlayer{track: Track{
			ID: "1", Artist: "A", Title: "T", ArtURL: "file:///cover.jpg", Playing: true,
		}},
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{TelegramMinInterval: time.Second, DefaultFirstName: "Name"},
		Clock:    fakeClock{now: time.Unix(10, 0)},
		Emoji:    fixedEmoji{},
	})

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(tg.names) != 1 || len(tg.bios) != 1 || tg.avatars != 1 || tg.publicRemovals != 1 {
		t.Fatalf("updates = %d names, %d bios, %d avatars, %d public removals", len(tg.names), len(tg.bios), tg.avatars, tg.publicRemovals)
	}
	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(tg.bios) != 1 || tg.avatars != 1 {
		t.Fatalf("duplicate updates = %d bios, %d avatars", len(tg.bios), tg.avatars)
	}
}

func TestTickTreatsMissingPlayerAsIdle(t *testing.T) {
	tg := &fakeTelegram{removeLimit: 1}
	service := New(Dependencies{
		Player:   fakePlayer{err: ErrNoPlayer},
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{
			DefaultAvatarPath: "default.jpg",
			DefaultFirstName:  "Name",
			IdleBio:           "Default bio",
		},
		Clock: fakeClock{now: time.Unix(10, 0)},
		Emoji: fixedEmoji{},
	})

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(tg.bios) != 1 || tg.bios[0] != "Default bio" {
		t.Fatalf("bios = %v", tg.bios)
	}
	if tg.avatars != 1 || tg.removals != 1 {
		t.Fatalf("avatars = %d removals = %d, want 1/1", tg.avatars, tg.removals)
	}
	if !service.noPlayerMode {
		t.Fatal("service should enter no-player mode")
	}
}

func TestTickReturnsTelegramError(t *testing.T) {
	expected := errors.New("telegram failed")
	service := New(Dependencies{
		Player:   fakePlayer{track: Track{Artist: "A", Title: "T", Playing: true}},
		Telegram: &fakeTelegram{err: expected},
		Avatar:   fakeAvatar{},
		Settings: Settings{TelegramMinInterval: time.Second, DefaultFirstName: "Name"},
		Clock:    fakeClock{now: time.Unix(10, 0)},
		Emoji:    fixedEmoji{},
	})

	if err := service.tick(context.Background()); !errors.Is(err, expected) {
		t.Fatalf("err = %v, want %v", err, expected)
	}
}
