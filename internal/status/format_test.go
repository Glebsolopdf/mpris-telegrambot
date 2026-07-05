package status

import "testing"

func TestMusicBioTruncatesToTelegramLimit(t *testing.T) {
	track := Track{
		Artist: "Очень Длинный Исполнитель Очень Длинный Исполнитель",
		Title:  "Очень Длинный Трек Очень Длинный Трек Очень Длинный Трек Очень Длинный Трек",
	}

	bio := MusicBio(track, "🎵", "")
	if len([]rune(bio)) > telegramBioLimit {
		t.Fatalf("bio has %d runes, want <= %d", len([]rune(bio)), telegramBioLimit)
	}
	if got := []rune(bio)[len([]rune(bio))-1]; got != '…' {
		t.Fatalf("last rune = %q, want ellipsis", got)
	}
}

func TestMusicBioUsesTemplateVariables(t *testing.T) {
	track := Track{Artist: "Artist", Title: "Track"}
	bio := MusicBio(track, "🎧", "{emoji_bio} Слушаю: {artist} — {title}")

	if bio != "🎧 Слушаю: Artist — Track" {
		t.Fatalf("bio = %q", bio)
	}
}

func TestMusicFirstNameUsesTemplateVariables(t *testing.T) {
	name := MusicFirstName("Substitute", "🎶", "{default_first_name} {emoji_name}")

	if name != "Substitute 🎶" {
		t.Fatalf("name = %q", name)
	}
}
