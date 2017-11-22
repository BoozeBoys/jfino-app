package tasks_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/BoozeBoys/jfino-app/tasks"
	"github.com/BoozeBoys/jfino-app/testutils"
)

func TestSendCommands(t *testing.T) {
	s := new(tasks.State)
	s.Power = true
	s.Motors[0].Speed = 255
	s.Motors[1].Speed = 255

	mc := new(testutils.CommanderMock)
	task := tasks.NewSendCommands(mc)

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	if len(mc.Commands) != 2 {
		t.FailNow()
	}

	t.Run("power", func(t *testing.T) {
		t.Parallel()

		want := fmt.Sprintf("%v", testutils.CommandDescription{Name: "power", Args: []interface{}{true}})
		got := fmt.Sprintf("%v", mc.Commands[0])

		if want != got {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})

	t.Run("speed", func(t *testing.T) {
		t.Parallel()

		want := fmt.Sprintf("%v", testutils.CommandDescription{Name: "speed", Args: []interface{}{255, 255}})
		got := fmt.Sprintf("%v", mc.Commands[1])

		if want != got {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
}

func TestSendCommandsError(t *testing.T) {
	s := new(tasks.State)

	mc := new(testutils.CommanderMock)
	task := tasks.NewSendCommands(mc)

	t.Run("power", func(t *testing.T) {
		t.Parallel()

		mc.SetError("power", errors.New("error"))
		mc.SetError("speed", nil)

		err := task.Perform(s)
		if err == nil {
			t.Errorf("want %v, got: %v", errors.New("error"), err)
		}
	})

	t.Run("speed", func(t *testing.T) {
		t.Parallel()

		mc.SetError("power", nil)
		mc.SetError("speed", errors.New("error"))

		err := task.Perform(s)
		if err == nil {
			t.Errorf("want %v, got: %v", errors.New("error"), err)
		}
	})
}
