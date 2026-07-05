package mpris

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestTrackFromProperties(t *testing.T) {
	track, err := TrackFromProperties(map[string]dbus.Variant{
		"PlaybackStatus": dbus.MakeVariant("Playing"),
		"Metadata": dbus.MakeVariant(map[string]dbus.Variant{
			"xesam:artist": dbus.MakeVariant([]string{"Artist 1", "Artist 2"}),
			"xesam:title":  dbus.MakeVariant("Track"),
			"mpris:artUrl": dbus.MakeVariant("file:///tmp/cover.jpg"),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !track.Playing || track.Artist != "Artist 1, Artist 2" || track.Title != "Track" {
		t.Fatalf("unexpected track: %+v", track)
	}
	if track.ID == "" {
		t.Fatal("fallback ID is empty")
	}
}
