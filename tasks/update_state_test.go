package tasks_test

import (
	"errors"
	"testing"

	"github.com/BoozeBoys/jfino-app/state"
	"github.com/BoozeBoys/jfino-app/tasks"
	"github.com/BoozeBoys/jfino-app/testutils"
)

func TestUpdateStatePower(t *testing.T) {
	r := new(testutils.StatusReaderMock)
	r.SetStatus([][]byte{
		[]byte("POWER 1"),
	})

	s := new(state.State)
	task := tasks.NewUpdateStatus(r)

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	if !s.Power {
		t.Errorf("want: %v, got: %v", true, s.Power)
	}
}

func TestUpdateStateSpeed(t *testing.T) {
	r := new(testutils.StatusReaderMock)
	r.SetStatus([][]byte{
		[]byte("SPEED 0 255"),
		[]byte("SPEED 1 -255"),
	})

	s := new(state.State)
	task := tasks.NewUpdateStatus(r)

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	if s.Motors[0].Speed != 255 {
		t.Errorf("want: %v, got: %v", 255, s.Motors[0].Speed)
	}

	if s.Motors[1].Speed != -255 {
		t.Errorf("want: %v, got: %v", -255, s.Motors[1].Speed)
	}
}

func TestUpdateStateCurrent(t *testing.T) {
	r := new(testutils.StatusReaderMock)
	r.SetStatus([][]byte{
		[]byte("CURRENT 0 1023"),
		[]byte("CURRENT 1 1023"),
	})

	s := new(state.State)
	task := tasks.NewUpdateStatus(r)

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	if s.Motors[0].Current != 1023 {
		t.Errorf("want: %v, got: %v", 1023, s.Motors[0].Current)
	}

	if s.Motors[1].Current != 1023 {
		t.Errorf("want: %v, got: %v", 1023, s.Motors[1].Current)
	}
}

func TestUpdateAnchorReport(t *testing.T) {
	r := new(testutils.StatusReaderMock)
	r.SetStatus([][]byte{
		[]byte("ANCHOR 1 12.45 -90"),
		[]byte("ANCHOR 3 1.66 -75"),
	})

	s := new(state.State)
	task := tasks.NewUpdateStatus(r)

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	check := state.AnchorReport{Range: 12.45, Power: -90}
	if s.RangeReport[1] != check {
		t.Errorf("want: %v, got: %v", check, s.RangeReport)
	}

	check = state.AnchorReport{Range: 1.66, Power: -75}
	if s.RangeReport[3] != check {
		t.Errorf("want: %v, got: %v", check, s.RangeReport)
	}
}

func TestUpdateStateError(t *testing.T) {
	r := new(testutils.StatusReaderMock)
	r.SetError(errors.New("error"))

	s := new(state.State)
	task := tasks.NewUpdateStatus(r)

	if err := task.Perform(s); err == nil {
		t.Errorf("want: %v got: %v", errors.New("error"), err)
	}
}
