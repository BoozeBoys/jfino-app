package tasks

import "github.com/BoozeBoys/jfino-app/state"
import "github.com/BoozeBoys/jfino-app/loc"

type Navigate struct {
	threshold loc.Meters
}

func NewNavigate(threshold loc.Meters) *Navigate {
	return &Navigate{threshold: threshold}
}

func (c *Navigate) Perform(s *state.State) error {
	if !s.Autopilot {
		return nil
	}

	s.Course = s.DestinationPoint.Sub(s.CurrentPosition)
	cp := s.CurrentPosition
	dp := s.DestinationPoint
	cp[2] = 0
	dp[2] = 0
	if cp.Distance(dp) < c.threshold {
		s.Autopilot = false
	}
	return nil
}