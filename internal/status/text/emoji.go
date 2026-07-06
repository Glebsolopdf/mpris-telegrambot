package text

import (
	"math/rand"
	"strings"
)

type EmojiPicker interface {
	Pick() string
}

type RandomEmojiPicker struct{}

func (RandomEmojiPicker) Pick() string {
	return musicEmojis[rand.Intn(len(musicEmojis))]
}

func ActiveEmoji(setting string, fallback string) string {
	value := strings.TrimSpace(setting)
	if value == "" {
		return fallback
	}
	if isDisabledEmoji(value) {
		return ""
	}
	options := emojiOptions(value)
	if len(options) == 0 {
		return fallback
	}
	return options[rand.Intn(len(options))]
}

func isDisabledEmoji(value string) bool {
	switch strings.ToLower(value) {
	case "0", "false", "no", "off", "disabled":
		return true
	default:
		return false
	}
}

func emojiOptions(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	options := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			options = append(options, trimmed)
		}
	}
	return options
}
