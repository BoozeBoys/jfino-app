package tasks

type Commander interface {
	Power(bool) error
	Speed(int, int) error
}

type SendCommands struct {
	c Commander
}

func NewSendCommands(c Commander) *SendCommands {
	return &SendCommands{c: c}
}

func (c *SendCommands) Perform(s *State) error {
	if err := c.c.Power(s.Power); err != nil {
		return err
	}

	if err := c.c.Speed(s.Motors[0].Speed, s.Motors[1].Speed); err != nil {
		return err
	}

	return nil
}
