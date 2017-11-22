package tasks

import (
	"bytes"
	"strconv"
)

type StatusReader interface {
	Status() ([][]byte, error)
}

type UpdateState struct {
	r StatusReader
}

func NewUpdateStatus(r StatusReader) *UpdateState {
	return &UpdateState{r: r}
}

func (c *UpdateState) Perform(s *State) error {
	lines, err := c.r.Status()
	if err != nil {
		return err
	}

	for _, line := range lines {
		tokens := bytes.Split(line, []byte{' '})
		switch string(tokens[0]) {
		case "POWER":
			s.Power = bytes.Equal(tokens[1], []byte{'1'})
		case "SPEED":
			motorID, err := strconv.Atoi(string(tokens[1]))
			if err != nil {
				break
			}

			speed, err := strconv.Atoi(string(tokens[2]))
			if err != nil {
				break
			}

			s.Motors[motorID].Speed = speed
		case "CURRENT":
			motorID, err := strconv.Atoi(string(tokens[1]))
			if err != nil {
				break
			}

			current, err := strconv.Atoi(string(tokens[2]))
			if err != nil {
				break
			}

			s.Motors[motorID].Current = current
		}
	}

	return nil
}
