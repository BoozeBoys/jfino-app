package tasks_test

import (
	"fmt"
	"testing"

	"github.com/BoozeBoys/jfino-app/commander"
	"github.com/BoozeBoys/jfino-app/state"
	"github.com/BoozeBoys/jfino-app/tasks"
	"github.com/BoozeBoys/jfino-app/testutils"
)

func TestWaitHelo(t *testing.T) {
	cases := []struct {
		connected bool
		err       error
		expected  bool
	}{
		{false, nil, true},
		{true, nil, true},
		{false, commander.Error, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("connected=%v error=%v", tc.connected, tc.err), func(tt *testing.T) {
			tt.Parallel()

			s := &state.State{ConnectedToFW: tc.connected}
			h := &testutils.HeloWaiterMock{Err: tc.err}
			task := tasks.NewWaitHelo(h)

			err := task.Perform(s)
			if err != tc.err {
				t.Errorf("want: %v, got: %v", tc.err, err)
			}

			if s.ConnectedToFW != tc.expected {
				t.Errorf("want: %v, got: %v", tc.expected, s.ConnectedToFW)
			}
		})
	}
}
