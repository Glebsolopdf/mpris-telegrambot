package mpris

import (
	"context"
	"errors"
	"strings"

	"github.com/godbus/dbus/v5"

	"mpris-tg-status/internal/status"
)

const playerPath = dbus.ObjectPath("/org/mpris/MediaPlayer2")

type Client struct {
	conn      *dbus.Conn
	preferred []string
}

func NewClient(preferred string) (*Client, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, preferred: normalizedPreferred(preferred)}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Current(ctx context.Context) (status.Track, error) {
	names, err := c.names(ctx)
	if err != nil {
		return status.Track{}, err
	}

	if len(c.preferred) > 0 {
		return c.preferredTrack(ctx, names)
	}

	var fallback status.Track
	for _, name := range names {
		track, err := c.track(ctx, name)
		if err != nil {
			continue
		}
		if track.Playing {
			return track, nil
		}
		if fallback.ID == "" {
			fallback = track
		}
	}

	if fallback.ID != "" {
		return fallback, nil
	}
	return status.Track{}, status.ErrNoPlayer
}

func (c *Client) preferredTrack(ctx context.Context, names []string) (status.Track, error) {
	for _, name := range names {
		if !matchesPreferred(name, c.preferred) {
			continue
		}
		track, err := c.track(ctx, name)
		if err != nil {
			return status.Track{}, err
		}
		return track, nil
	}
	return status.Track{}, status.ErrNoPlayer
}

func (c *Client) names(ctx context.Context) ([]string, error) {
	var names []string
	call := c.conn.BusObject().GoWithContext(ctx, "org.freedesktop.DBus.ListNames", 0, nil)
	select {
	case done := <-call.Done:
		if done.Err != nil {
			return nil, done.Err
		}
		if err := done.Store(&names); err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	players := make([]string, 0)
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			players = append(players, name)
		}
	}
	return players, nil
}

func (c *Client) track(ctx context.Context, name string) (status.Track, error) {
	props, err := c.properties(ctx, name)
	if err != nil {
		return status.Track{}, err
	}
	return TrackFromProperties(props)
}

func (c *Client) properties(ctx context.Context, name string) (map[string]dbus.Variant, error) {
	var props map[string]dbus.Variant
	obj := c.conn.Object(name, playerPath)
	call := obj.GoWithContext(ctx, "org.freedesktop.DBus.Properties.GetAll", 0, nil, "org.mpris.MediaPlayer2.Player")
	select {
	case done := <-call.Done:
		if done.Err != nil {
			return nil, done.Err
		}
		if err := done.Store(&props); err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if len(props) == 0 {
		return nil, errors.New("empty MPRIS properties")
	}
	return props, nil
}
