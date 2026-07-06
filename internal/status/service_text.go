package status

import (
	"context"

	statustext "mpris-tg-status/internal/status/text"
)

func (s *Service) applyActiveText(ctx context.Context, track Track) error {
	now := s.deps.Clock.Now()
	firstName := statustext.MusicFirstNameAt(
		s.deps.Settings.DefaultFirstName,
		s.activeNameEmoji,
		s.deps.Settings.ActiveFirstNameTemplate,
		now,
	)
	if err := s.applyName(ctx, firstName, s.deps.Settings.DefaultLastName); err != nil {
		return err
	}

	bio := statustext.MusicBioAt(textTrack(track), s.activeBioEmoji, s.deps.Settings.ActiveBioTemplate, now)
	return s.applyBio(ctx, bio)
}

func textTrack(track Track) statustext.Track {
	return statustext.Track{
		Artist: track.Artist,
		Title:  track.Title,
		URL:    track.URL,
	}
}
