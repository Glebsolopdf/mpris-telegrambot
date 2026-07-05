package status

import "testing"

func TestActiveEmojiUsesFallbackWhenUnset(t *testing.T) {
	if got := ActiveEmoji("", "🎵"); got != "🎵" {
		t.Fatalf("emoji = %q", got)
	}
}

func TestActiveEmojiCanBeDisabled(t *testing.T) {
	if got := ActiveEmoji("false", "🎵"); got != "" {
		t.Fatalf("emoji = %q, want empty", got)
	}
}

func TestActiveEmojiUsesConfiguredSingleValue(t *testing.T) {
	if got := ActiveEmoji("🔥", "🎵"); got != "🔥" {
		t.Fatalf("emoji = %q", got)
	}
}
