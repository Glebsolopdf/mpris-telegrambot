package status

import "strings"

const telegramBioLimit = 140

const (
	defaultActiveBioTemplate       = "{emoji_bio} Listening now: {artist} -- {title}"
	defaultActiveFirstNameTemplate = "{emoji_name} {default_first_name}"
)

var musicEmojis = []string{
	"🎵", "🎶", "〽️", "📻", "🎧", "🎸", "🎹", "🪗", "🥁",
	"🪘", "🎺", "🎻", "🪉", "🪕", "🪈", "🪇", "🔊", "🎼",
}

func MusicBio(track Track, emoji string, template string) string {
	if template == "" {
		template = defaultActiveBioTemplate
	}
	return truncateRunes(renderTemplate(template, templateVars{
		Emoji:    emoji,
		EmojiBio: emoji,
		Artist:   track.Artist,
		Title:    track.Title,
	}), telegramBioLimit)
}

func MusicFirstName(defaultName string, emoji string, template string) string {
	if template == "" {
		template = defaultActiveFirstNameTemplate
	}
	return truncateRunes(renderTemplate(template, templateVars{
		Emoji:            emoji,
		EmojiName:        emoji,
		DefaultFirstName: defaultName,
	}), 64)
}

type templateVars struct {
	Emoji            string
	EmojiName        string
	EmojiBio         string
	Artist           string
	Title            string
	DefaultFirstName string
}

func renderTemplate(template string, vars templateVars) string {
	replacer := strings.NewReplacer(
		"{emoji}", vars.Emoji,
		"{emoji_name}", vars.EmojiName,
		"{emoji_bio}", vars.EmojiBio,
		"{artist}", vars.Artist,
		"{title}", vars.Title,
		"{default_first_name}", vars.DefaultFirstName,
	)
	return strings.TrimSpace(replacer.Replace(template))
}

func truncateRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "…"
}
