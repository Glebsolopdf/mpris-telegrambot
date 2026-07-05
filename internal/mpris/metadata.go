package mpris

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"

	"mpris-tg-status/internal/status"
)

func TrackFromProperties(props map[string]dbus.Variant) (status.Track, error) {
	metadata, ok := variantMap(props["Metadata"])
	if !ok {
		return status.Track{}, fmt.Errorf("missing MPRIS metadata")
	}

	track := status.Track{
		ID:      stringValue(metadata["mpris:trackid"]),
		Artist:  artistsValue(metadata["xesam:artist"]),
		Title:   stringValue(metadata["xesam:title"]),
		ArtURL:  stringValue(metadata["mpris:artUrl"]),
		Playing: stringValue(props["PlaybackStatus"]) == "Playing",
	}
	if track.ID == "" {
		track.ID = track.Artist + "\x00" + track.Title + "\x00" + track.ArtURL
	}
	return track, nil
}

func variantMap(value dbus.Variant) (map[string]dbus.Variant, bool) {
	metadata, ok := value.Value().(map[string]dbus.Variant)
	return metadata, ok
}

func stringValue(value dbus.Variant) string {
	switch typed := value.Value().(type) {
	case string:
		return strings.TrimSpace(typed)
	case dbus.ObjectPath:
		return string(typed)
	default:
		return ""
	}
}

func artistsValue(value dbus.Variant) string {
	switch typed := value.Value().(type) {
	case []string:
		return strings.Join(typed, ", ")
	case string:
		return typed
	default:
		return ""
	}
}
