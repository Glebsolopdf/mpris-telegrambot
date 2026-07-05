package status

import (
	"context"
	"testing"
	"time"
)

func TestRunReturnsBusinessConnectionInvalid(t *testing.T) {
	expected := fakeBusinessInvalid{}
	service := New(Dependencies{
		Player: fakePlayer{track: Track{
			ID: "1", Artist: "A", Title: "T", Playing: true,
		}},
		Telegram: &fakeTelegram{err: expected},
		Settings: Settings{
			DefaultFirstName:    "Name",
			TelegramMinInterval: time.Second,
		},
		Clock: fakeClock{now: time.Unix(10, 0)},
		Emoji: fixedEmoji{},
	})

	if err := service.Run(context.Background()); err != expected {
		t.Fatalf("err = %v, want invalid", err)
	}
}

type fakeBusinessInvalid struct{}

func (fakeBusinessInvalid) Error() string { return "business invalid" }

func (fakeBusinessInvalid) BusinessConnectionInvalid() bool { return true }
