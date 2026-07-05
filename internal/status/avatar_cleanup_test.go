package status

import (
	"context"
	"testing"
	"time"
)

func TestTickSwitchesAlbumBackToDefaultAvatar(t *testing.T) {
	tg := &fakeTelegram{}
	service := newAvatarTestService(tg)

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	service.deps.Player = fakePlayer{track: Track{Playing: false}}
	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	if tg.removals != 2 || tg.avatars != 2 {
		t.Fatalf("removals = %d avatars = %d, want 2/2", tg.removals, tg.avatars)
	}
}

func TestTickRemovesUnknownMainAvatarForSinglePhotoMode(t *testing.T) {
	tg := &fakeTelegram{removeLimit: 1}
	service := newAvatarTestService(tg)

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	if tg.removals != 1 || tg.avatars != 1 {
		t.Fatalf("removals = %d avatars = %d, want 1/1", tg.removals, tg.avatars)
	}
}

func TestAvatarCooldownSkipsAvatarButKeepsTextUpdates(t *testing.T) {
	tg := &fakeTelegram{}
	service := newAvatarTestService(tg)
	service.nextAvatar = service.deps.Clock.Now().Add(time.Hour)

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(tg.names) != 1 || len(tg.bios) != 1 {
		t.Fatalf("text updates = %d names, %d bios", len(tg.names), len(tg.bios))
	}
	if tg.avatars != 0 || tg.removals != 0 || tg.publicRemovals != 0 {
		t.Fatalf("avatar updates = %d/%d/%d, want none", tg.avatars, tg.removals, tg.publicRemovals)
	}
}

func TestGeneratedAvatarDisabledKeepsTextUpdates(t *testing.T) {
	tg := &fakeTelegram{}
	service := newAvatarTestService(tg)
	service.deps.Settings.DisableGeneratedAvatar = true

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(tg.names) != 1 || len(tg.bios) != 1 {
		t.Fatalf("text updates = %d names, %d bios", len(tg.names), len(tg.bios))
	}
	if tg.avatars != 0 || tg.removals != 0 || tg.publicRemovals != 0 {
		t.Fatalf("avatar updates = %d/%d/%d, want none", tg.avatars, tg.removals, tg.publicRemovals)
	}
}

func TestSameCoverDifferentTrackSkipsAvatarUpdate(t *testing.T) {
	tg := &fakeTelegram{}
	service := newAvatarTestService(tg)

	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}
	service.nextAvatar = time.Time{}
	service.deps.Player = fakePlayer{track: Track{
		ID: "2", Artist: "A", Title: "Second", ArtURL: "file:///cover.jpg", Playing: true,
	}}
	if err := service.tick(context.Background()); err != nil {
		t.Fatal(err)
	}

	if tg.avatars != 1 || tg.removals != 1 {
		t.Fatalf("avatar updates = %d removals = %d, want 1/1", tg.avatars, tg.removals)
	}
}

func newAvatarTestService(tg *fakeTelegram) *Service {
	return New(Dependencies{
		Player: fakePlayer{track: Track{
			ID: "1", Artist: "A", Title: "T", ArtURL: "file:///cover.jpg", Playing: true,
		}},
		Telegram: tg,
		Avatar:   fakeAvatar{},
		Settings: Settings{DefaultAvatarPath: "default.jpg", DefaultFirstName: "Name"},
		Clock:    fakeClock{now: time.Unix(10, 0)},
		Emoji:    fixedEmoji{},
	})
}
