package tasks_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/BoozeBoys/jfino-app/tasks"

	"github.com/BoozeBoys/jfino-app/loc"
	"github.com/BoozeBoys/jfino-app/state"
)

func TestErrorRms(t *testing.T) {
	ranges := make(map[string]state.AnchorReport)
	ranges["0"] = state.AnchorReport{Range: 1}
	ranges["1"] = state.AnchorReport{Range: 1}
	ranges["2"] = state.AnchorReport{Range: 1}
	ranges["3"] = state.AnchorReport{Range: 1}

	d := loc.Meters(1.0)
	anchors := make(map[string]tasks.AnchorCfg)
	assign := func() {
		anchors["0"] = tasks.AnchorCfg{Loc: loc.Point{-d, 0, 0}, Offset: 0}
		anchors["1"] = tasks.AnchorCfg{Loc: loc.Point{0, d, 0}, Offset: 0}
		anchors["2"] = tasks.AnchorCfg{Loc: loc.Point{d, 0, 0}, Offset: 0}
		anchors["3"] = tasks.AnchorCfg{Loc: loc.Point{0, -d, 0}, Offset: 0}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	rms := ep.ErrorPtoP(ranges, loc.Point{0, 0})
	if math.Abs(float64(rms)-0) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	rms = ep.ErrorPtoP(ranges, loc.Point{1, 1})
	res := math.Sqrt((math.Pow(math.Sqrt(1+4)-1, 2)*2)/4) * 3
	if math.Abs(float64(rms)-res) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	d = 3.0
	assign()
	rms = ep.ErrorPtoP(ranges, loc.Point{0, 0})
	if math.Abs(float64(rms)-6) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}
}

func TestFindBoundingBox(t *testing.T) {
	ranges := make(map[string]state.AnchorReport)
	ranges["0"] = state.AnchorReport{Range: 10}
	ranges["1"] = state.AnchorReport{Range: 10}
	ranges["2"] = state.AnchorReport{Range: 10}
	ranges["3"] = state.AnchorReport{Range: 10}

	d := loc.Meters(1.0)
	anchors := make(map[string]tasks.AnchorCfg)
	assign := func() {
		anchors["0"] = tasks.AnchorCfg{Loc: loc.Point{-d, 0, 1}, Offset: 0}
		anchors["1"] = tasks.AnchorCfg{Loc: loc.Point{0, d, 1}, Offset: 0}
		anchors["2"] = tasks.AnchorCfg{Loc: loc.Point{d, 0, 1}, Offset: 0}
		anchors["3"] = tasks.AnchorCfg{Loc: loc.Point{0, -d, 1}, Offset: 0}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	b := ep.FindBoundingBox(ranges)

	if !b[0].IsEqual(loc.Point{-11, -11, -9}) {
		t.Fatalf("bo %v", b)
	}

	if !b[1].IsEqual(loc.Point{11, 11, 11}) {
		t.Fatalf("bo %v", b)
	}
}

func TestComputePosition(t *testing.T) {
	anchors := make(map[string]tasks.AnchorCfg)

	anchors["0"] = tasks.AnchorCfg{Loc: loc.Point{0, 0, 0}, Offset: 0}
	anchors["1"] = tasks.AnchorCfg{Loc: loc.Point{0, 100, 0}, Offset: 0}
	anchors["2"] = tasks.AnchorCfg{Loc: loc.Point{100, 0, 0}, Offset: 0}
	anchors["3"] = tasks.AnchorCfg{Loc: loc.Point{0, 0, 100}, Offset: 0}

	r := rand.New(rand.NewSource(0))
	ep := tasks.NewEstimatePosition(anchors)
	errCount := 0
	cnt := 100000
	for i := 0; i < cnt; i++ {
		ranges := make(map[string]state.AnchorReport)
		x := loc.Meters(r.Float64() * 100)
		y := loc.Meters(r.Float64() * 100)
		z := loc.Meters(r.Float64() * 100)
		j := loc.Point{x, y, z}

		err := loc.Meters(0.2)
		for i, a := range anchors {
			ranges[i] = state.AnchorReport{Range: j.Distance(a.Loc) + a.Offset + loc.Meters(r.Float64())*err - err/2}
		}

		p, _ := ep.ComputePosition(ep.FindBoundingBox(ranges), ranges)

		//fmt.Printf("p %v, accuracy +/-%.2f, actual dist %f, \n", p, float64(err), p.Distance(j))
		if p.Distance(j) > err {
			errCount++
			fmt.Println("Error:", errCount)
			fmt.Printf("p %v, accuracy +/-%.2f, actual dist %f, \n", p, float64(err), p.Distance(j))
		}
	}
	if errCount > int(float64(cnt)*(1-0.99)) {
		t.Fatalf("want: %v, got %v", float64(cnt)*(1-0.99), errCount)
	}
}

func TestEstimatePosition(t *testing.T) {
	anchors := make(map[string]tasks.AnchorCfg)

	anchors["0"] = tasks.AnchorCfg{Loc: loc.Point{0, 0, 0}, Offset: 0}
	anchors["1"] = tasks.AnchorCfg{Loc: loc.Point{0, 100, 0}, Offset: 0}
	anchors["2"] = tasks.AnchorCfg{Loc: loc.Point{100, 0, 0}, Offset: 0}
	anchors["3"] = tasks.AnchorCfg{Loc: loc.Point{0, 0, 100}, Offset: 0}
	s := new(state.State)
	task := tasks.NewEstimatePosition(anchors)

	ranges := make(map[string]state.AnchorReport)
	j := loc.Point{10, -34, 23}

	for i, a := range anchors {
		ranges[i] = state.AnchorReport{Range: j.Distance(a.Loc) + a.Offset}
	}
	s.RangeReport = ranges

	if err := task.Perform(s); err != nil {
		t.FailNow()
	}

	if s.PositionAccuracy > 0.005 {
		t.Errorf("want < 0.005, got %v", s.PositionAccuracy)
	}

	if s.CurrentPosition.Distance(j) > 0.01 {
		t.Errorf("want %v, got %v", j, s.CurrentPosition)
	}
}
