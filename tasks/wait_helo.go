package tasks

import "github.com/BoozeBoys/jfino-app/state"

type HeloWaiter interface {
	WaitHelo() error
}

type WaitHelo struct {
	r HeloWaiter
}

func NewWaitHelo(r HeloWaiter) *WaitHelo {
	return &WaitHelo{r: r}
}

func (t *WaitHelo) Perform(s *state.State) error {
	if s.ConnectedToFW {
		return nil
	}

	if err := t.r.WaitHelo(); err != nil {
		s.ConnectedToFW = false
		return err
	}

	s.ConnectedToFW = true

	return nil
}
