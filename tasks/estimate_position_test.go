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
		anchors[0] = loc.Point{-d, 0, 0}
		anchors[1] = loc.Point{0, d, 0}
		anchors[2] = loc.Point{d, 0, 0}
		anchors[3] = loc.Point{0, -d, 0}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	rms := ep.ErrorRms(ranges, loc.Point{0, 0})
	if math.Abs(rms-0) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	rms = ep.ErrorRms(ranges, loc.Point{1, 1})
	res := math.Sqrt((math.Pow(math.Sqrt(1+4)-1, 2) * 2) / 4)
	if math.Abs(rms-res) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	d = 3.0
	assign()
	rms = ep.ErrorRms(ranges, loc.Point{0, 0})
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
		anchors[0] = loc.Point{-d, 0, 1}
		anchors[1] = loc.Point{0, d, 1}
		anchors[2] = loc.Point{d, 0, 1}
		anchors[3] = loc.Point{0, -d, 1}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	b := ep.FindBoundingBox(ranges)

	if !b.P0.IsEqual(loc.Point{-11, -11, -9}) {
		t.Fatalf("bo %v", b)
	}

	if !b.P1.IsEqual(loc.Point{11, 11, 11}) {
		t.Fatalf("bo %v", b)
	}
}
