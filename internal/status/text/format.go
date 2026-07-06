package text

import "time"

const telegramBioLimit = 140

const (
	defaultActiveBioTemplate       = "{emoji_bio} {time} Listening now: {artist} -- {title}"
	defaultActiveFirstNameTemplate = "{emoji_name} {default_first_name}"
)

var musicEmojis = []string{
	"🎵", "🎶", "〽️", "📻", "🎧", "🎸", "🎹", "🪗", "🥁",
	"🪘", "🎺", "🎻", "🪉", "🪕", "🪈", "🪇", "🔊", "🎼",
}

type Track struct {
	Artist string
	Title  string
	URL    string
}

func MusicBio(track Track, emoji string, template string) string {
	return MusicBioAt(track, emoji, template, time.Now())
}

func MusicBioAt(track Track, emoji string, template string, now time.Time) string {
	if template == "" {
		template = defaultActiveBioTemplate
	}
	return fitBioTemplate(template, templateVars{
		Emoji:    emoji,
		EmojiBio: emoji,
		Artist:   track.Artist,
		Title:    track.Title,
		Time:     formatTemplateTime(now),
		TrackURL: track.URL,
	})
}

func MusicFirstName(defaultName string, emoji string, template string) string {
	return MusicFirstNameAt(defaultName, emoji, template, time.Now())
}

func MusicFirstNameAt(defaultName string, emoji string, template string, now time.Time) string {
	if template == "" {
		template = defaultActiveFirstNameTemplate
	}
	return truncateRunes(renderTemplate(template, templateVars{
		Emoji:            emoji,
		EmojiName:        emoji,
		DefaultFirstName: defaultName,
		Time:             formatTemplateTime(now),
	}), 64)
}

type templateVars struct {
	Emoji            string
	EmojiName        string
	EmojiBio         string
	Artist           string
	Title            string
	Time             string
	TrackURL         string
	DefaultFirstName string
}

func formatTemplateTime(now time.Time) string {
	return now.Format("15:04")
}
