package text

import "strings"

func renderTemplate(template string, vars templateVars) string {
	template = strings.ReplaceAll(template, `\n`, "\n")
	replacer := strings.NewReplacer(
		"{emoji}", vars.Emoji,
		"{emoji_name}", vars.EmojiName,
		"{emoji_bio}", vars.EmojiBio,
		"{artist}", vars.Artist,
		"{title}", vars.Title,
		"{time}", vars.Time,
		"{track_url}", vars.TrackURL,
		"{url}", vars.TrackURL,
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
