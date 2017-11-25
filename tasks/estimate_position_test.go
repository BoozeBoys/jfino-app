package tasks_test

import (
	"math"
	"testing"

	"github.com/BoozeBoys/jfino-app/tasks"

	"github.com/BoozeBoys/jfino-app/loc"
	"github.com/BoozeBoys/jfino-app/state"
)

func TestErrorRms(t *testing.T) {
	ranges := make(map[uint]state.AnchorReport)
	ranges[0] = state.AnchorReport{Range: 1}
	ranges[1] = state.AnchorReport{Range: 1}
	ranges[2] = state.AnchorReport{Range: 1}
	ranges[3] = state.AnchorReport{Range: 1}

	d := 1.0
	anchors := make(map[uint]loc.Point)
	assign := func() {
		anchors[0] = loc.Point{X: -d, Y: 0, Z: 0}
		anchors[1] = loc.Point{X: 0, Y: d, Z: 0}
		anchors[2] = loc.Point{X: d, Y: 0, Z: 0}
		anchors[3] = loc.Point{X: 0, Y: -d, Z: 0}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	rms := ep.ErrorRms(ranges, loc.Point{X: 0, Y: 0})
	if math.Abs(rms-0) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	rms = ep.ErrorRms(ranges, loc.Point{X: 1, Y: 1})
	res := math.Sqrt((math.Pow(math.Sqrt(1+4)-1, 2) * 2) / 4)
	if math.Abs(rms-res) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	d = 3.0
	assign()
	rms = ep.ErrorRms(ranges, loc.Point{X: 0, Y: 0})
	if math.Abs(rms-2) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}
}

func TestFindBoundingBox(t *testing.T) {
	ranges := make(map[uint]state.AnchorReport)
	ranges[0] = state.AnchorReport{Range: 10}
	ranges[1] = state.AnchorReport{Range: 10}
	ranges[2] = state.AnchorReport{Range: 10}
	ranges[3] = state.AnchorReport{Range: 10}

	d := 1.0
	anchors := make(map[uint]loc.Point)
	assign := func() {
		anchors[0] = loc.Point{X: -d, Y: 0, Z: 1}
		anchors[1] = loc.Point{X: 0, Y: d, Z: 1}
		anchors[2] = loc.Point{X: d, Y: 0, Z: 1}
		anchors[3] = loc.Point{X: 0, Y: -d, Z: 1}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	b := ep.FindBoundingBox(ranges)

	if !b.P0.IsEqual(loc.Point{X: -11, Y: -11, Z: -9}) {
		t.Fatalf("box: %v", b)
	}

	if !b.P1.IsEqual(loc.Point{X: 11, Y: 11, Z: 11}) {
		t.Fatalf("box: %v", b)
	}
}
