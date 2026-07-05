package status

import (
	"context"
	"errors"
	"time"
)

type fakeClock struct{ now time.Time }

func (c fakeClock) Now() time.Time { return c.now }

type fakePlayer struct {
	track Track
	err   error
}

func (p fakePlayer) Current(context.Context) (Track, error) {
	return p.track, p.err
}

type fakeTelegram struct {
	bios           []string
	names          []string
	avatars        int
	removals       int
	publicRemovals int
	removeLimit    int
	err            error
}

func (t *fakeTelegram) SetBio(_ context.Context, bio string) error {
	t.bios = append(t.bios, bio)
	return t.err
}

func (t *fakeTelegram) SetName(_ context.Context, firstName string, lastName string) error {
	t.names = append(t.names, firstName+" "+lastName)
	return t.err
}

func (t *fakeTelegram) SetAvatar(context.Context, []byte) error {
	t.avatars++
	return t.err
}

func (t *fakeTelegram) RemoveAvatar(context.Context) error {
	if t.removeLimit > 0 && t.removals >= t.removeLimit {
		return errors.New("no profile photo")
	}
	t.removals++
	return t.err
}

func (t *fakeTelegram) RemovePublicAvatar(context.Context) error {
	t.publicRemovals++
	return t.err
}

type fakeAvatar struct{}

func (fakeAvatar) Build(context.Context, string, string) ([]byte, error) {
	return []byte("jpg"), nil
}

func (fakeAvatar) LoadDefault(context.Context, string) ([]byte, error) {
	return []byte("default"), nil
}

type fixedEmoji struct{}

func (fixedEmoji) Pick() string { return "🎵" }
