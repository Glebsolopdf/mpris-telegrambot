package mpris

import "testing"

func TestMatchesPreferred(t *testing.T) {
	if !matchesPreferred("org.mpris.MediaPlayer2.spotify", "spotify") {
		t.Fatal("spotify player should match preferred name")
	}
	if matchesPreferred("org.mpris.MediaPlayer2.firefox", "spotify") {
		t.Fatal("firefox player should not match spotify preference")
	}
}

func TestNormalizedPreferredTreatsAllAsAnyPlayer(t *testing.T) {
	if got := normalizedPreferred("all"); got != "" {
		t.Fatalf("preferred = %q, want empty", got)
	}
	if got := normalizedPreferred(" ALL "); got != "" {
		t.Fatalf("preferred = %q, want empty", got)
	}
	if got := normalizedPreferred("spotify"); got != "spotify" {
		t.Fatalf("preferred = %q, want spotify", got)
	}
}
