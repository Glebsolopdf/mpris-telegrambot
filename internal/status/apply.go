package status

import (
	"context"
	"runtime/debug"
)

func (s *Service) applyBio(ctx context.Context, bio string) error {
	if bio == s.lastBio {
		return nil
	}
	if err := s.deps.Telegram.SetBio(ctx, bio); err != nil {
		s.pauseTextOnRetry(err)
		return err
	}
	s.log("bio updated: %q", bio)
	s.lastBio = bio
	s.nextText = s.deps.Clock.Now().Add(s.deps.Settings.TelegramMinInterval)
	return nil
}

func (s *Service) applyName(ctx context.Context, firstName string, lastName string) error {
	sameName := firstName == s.lastFirstName && lastName == s.lastLastName
	if firstName == "" || sameName {
		return nil
	}
	if err := s.deps.Telegram.SetName(ctx, firstName, lastName); err != nil {
		s.pauseTextOnRetry(err)
		return err
	}
	s.log("name updated: %q %q", firstName, lastName)
	s.lastFirstName = firstName
	s.lastLastName = lastName
	s.nextText = s.deps.Clock.Now().Add(s.deps.Settings.TelegramMinInterval)
	return nil
}

func (s *Service) applyAvatar(ctx context.Context, track Track) error {
	if s.deps.Settings.DisableGeneratedAvatar {
		s.debug("generated avatar update skipped by configuration")
		return nil
	}

	key := track.ArtURL
	sameAlbumAvatar := s.avatarMode == avatarAlbum && key == s.lastAvatarKey
	if track.ArtURL == "" || sameAlbumAvatar {
		return nil
	}
	if !s.avatarAllowed() {
		return nil
	}

	jpeg, err := s.deps.Avatar.Build(ctx, track.ArtURL, track.Title)
	if err != nil {
		return err
	}
	s.log("generated avatar for track: %q", track.Title)
	if err := s.switchAvatar(ctx, avatarAlbum, jpeg); err != nil {
		return err
	}
	debug.FreeOSMemory()
	s.lastAvatarKey = key
	return nil
}

func (s *Service) applyDefaultAvatar(ctx context.Context) error {
	if s.deps.Settings.DefaultAvatarPath == "" || s.avatarMode == avatarDefault {
		return nil
	}
	if !s.avatarAllowed() {
		return nil
	}
	image, err := s.deps.Avatar.LoadDefault(ctx, s.deps.Settings.DefaultAvatarPath)
	if err != nil {
		return err
	}
	s.log("loaded default avatar: %s", s.deps.Settings.DefaultAvatarPath)
	if err := s.switchAvatar(ctx, avatarDefault, image); err != nil {
		return err
	}
	debug.FreeOSMemory()
	s.lastAvatarKey = "default"
	return nil
}

func (s *Service) switchAvatar(ctx context.Context, mode avatarMode, image []byte) error {
	s.log("switching avatar to %s", mode)
	if err := s.removePublicAvatar(ctx); err != nil {
		return err
	}
	if err := s.removeMainAvatar(ctx); err != nil {
		return err
	}
	if err := s.deps.Telegram.SetAvatar(ctx, image); err != nil {
		s.pauseAvatarOnRetry(err)
		return err
	}
	s.avatarMode = mode
	s.log("%s avatar set as private main profile photo", mode)
	s.setAvatarCooldown(s.deps.Settings.AvatarMinInterval)
	return nil
}

func (s *Service) removePublicAvatar(ctx context.Context) error {
	if err := s.deps.Telegram.RemovePublicAvatar(ctx); err != nil {
		if retryDelay(err) > 0 {
			s.pauseAvatarOnRetry(err)
			return err
		}
		s.debug("public business avatar remove skipped: %v", err)
		return nil
	}
	s.debug("public business avatar removed")
	return nil
}
