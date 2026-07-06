package status

import "time"

func (s *Service) avatarAllowed() bool {
	if s.deps.Clock.Now().Before(s.nextAvatar) {
		s.debug("avatar cooldown active until %s", s.nextAvatar.Format(time.RFC3339))
		return false
	}
	return true
}

func (s *Service) setAvatarCooldown(delay time.Duration) {
	s.nextAvatar = s.deps.Clock.Now().Add(delay)
}

func (s *Service) pauseAvatarOnRetry(err error) {
	if delay := retryDelay(err); delay > 0 {
		s.setAvatarCooldown(delay)
		s.log("telegram avatar flood control: retry after %s", delay)
	}
}
