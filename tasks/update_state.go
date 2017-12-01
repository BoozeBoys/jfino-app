package tasks

import (
	"bytes"
	"strconv"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
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

func (c *UpdateState) Perform(s *state.State) error {
	lines, err := c.r.Status()
	if err != nil {
		return err
	}
	s.RangeReport = make(map[int]state.AnchorReport)

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
		case "ANCHOR":
			id, err := strconv.Atoi(string(tokens[1]))
			if err != nil {
				break
			}

			rng, err := strconv.ParseFloat(string(tokens[2]), 64)
			if err != nil {
				break
			}

			pwr, err := strconv.Atoi(string(tokens[3]))
			if err != nil {
				break
			}

			s.RangeReport[id] = state.AnchorReport{Range: loc.Meters(rng), Power: pwr}
		}
	}

	return nil
}
