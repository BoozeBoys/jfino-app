package tasks_test

import (
	"testing"
	"time"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/tasks"

	"github.com/BoozeBoys/jfino-app/state"
)

func TestHttpMap(t *testing.T) {
	s := new(state.State)
	task := tasks.NewHttpMap(8000, "../templates")

	for i := 0; i < 5; i++ {
		s.PositionAccuracy = loc.Meters(i)
		if err := task.Perform(s); err != nil {
			t.FailNow()
		}
		time.Sleep(1 * time.Second)
	}
}
