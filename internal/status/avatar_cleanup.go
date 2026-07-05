package status

import "context"

func (s *Service) removeMainAvatar(ctx context.Context) error {
	if err := s.deps.Telegram.RemoveAvatar(ctx); err != nil {
		if retryDelay(err) > 0 {
			s.pauseAvatarOnRetry(err)
			return err
		}
		s.debug("main business avatar remove skipped: %v", err)
		return nil
	}
	s.log("main business avatar removed")
	return nil
}
