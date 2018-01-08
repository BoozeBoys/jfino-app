package tasks_test

import (
	"fmt"
	"testing"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
	"github.com/BoozeBoys/jfino-app/tasks"
)

func TestNavigate(t *testing.T) {
	cases := []struct {
		CurrentPosition   loc.Point
		Destination       loc.Point
		Autopilot         bool
		ExpectedAutopilot bool
		ExpectedCourse    loc.Point
	}{
		{loc.Point{0, 0, 0}, loc.Point{2, 2, 2}, true, true, loc.Point{2, 2, 2}},
		{loc.Point{0, 0, 0}, loc.Point{0, 0, 0}, true, false, loc.Point{0, 0, 0}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("Autopilot=%v, Course=%v", tc.ExpectedAutopilot, tc.ExpectedCourse), func(tt *testing.T) {
			tt.Parallel()

			s := new(state.State)
			s.Autopilot = tc.Autopilot
			s.CurrentPosition = tc.CurrentPosition
			s.DestinationPoint = tc.Destination

			task := tasks.NewNavigate(1)

			err := task.Perform(s)
			if err != nil {
				tt.FailNow()
			}

			if s.Autopilot != tc.ExpectedAutopilot {
				tt.Errorf("got %v, want %v", s.Autopilot, tc.ExpectedAutopilot)
			}

			if !s.Course.IsEqual(tc.ExpectedCourse) {
				tt.Errorf("got %v, want %v", s.Course, tc.ExpectedCourse)
			}
		})
	}
}
