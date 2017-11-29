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
	ranges := make(map[int]state.AnchorReport)
	ranges[0] = state.AnchorReport{Range: 1}
	ranges[1] = state.AnchorReport{Range: 1}
	ranges[2] = state.AnchorReport{Range: 1}
	ranges[3] = state.AnchorReport{Range: 1}

	d := loc.Meters(1.0)
	anchors := make(map[int]loc.Point)
	assign := func() {
		anchors[0] = loc.Point{-d, 0, 0}
		anchors[1] = loc.Point{0, d, 0}
		anchors[2] = loc.Point{d, 0, 0}
		anchors[3] = loc.Point{0, -d, 0}
	}

	assign()
	ep := tasks.NewEstimatePosition(anchors)

	rms := ep.ErrorRms(ranges, loc.Point{0, 0})
	if math.Abs(float64(rms)-0) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	rms = ep.ErrorRms(ranges, loc.Point{1, 1})
	res := math.Sqrt((math.Pow(math.Sqrt(1+4)-1, 2) * 2) / 4)
	if math.Abs(float64(rms)-res) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}

	d = 3.0
	assign()
	rms = ep.ErrorRms(ranges, loc.Point{0, 0})
	if math.Abs(float64(rms)-2) > 1e-6 {
		t.Fatalf("rms: %f", rms)
	}
}

func TestFindBoundingBox(t *testing.T) {
	ranges := make(map[int]state.AnchorReport)
	ranges[0] = state.AnchorReport{Range: 10}
	ranges[1] = state.AnchorReport{Range: 10}
	ranges[2] = state.AnchorReport{Range: 10}
	ranges[3] = state.AnchorReport{Range: 10}

	d := loc.Meters(1.0)
	anchors := make(map[int]loc.Point)
	assign := func() {
		anchors[0] = loc.Point{-d, 0, 1}
		anchors[1] = loc.Point{0, d, 1}
		anchors[2] = loc.Point{d, 0, 1}
		anchors[3] = loc.Point{0, -d, 1}
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
	anchors := make(map[int]loc.Point)

	anchors[0] = loc.Point{0, 0, 0}
	anchors[1] = loc.Point{0, 100, 0}
	anchors[2] = loc.Point{100, 0, 0}
	anchors[3] = loc.Point{0, 0, 100}

	r := rand.New(rand.NewSource(0))
	ep := tasks.NewEstimatePosition(anchors)
	for i := 0; i < 10000; i++ {
		ranges := make(map[int]state.AnchorReport)
		x := loc.Meters(r.Float64() * 100)
		y := loc.Meters(r.Float64() * 100)
		z := loc.Meters(r.Float64() * 100)
		j := loc.Point{x, y, z}

		err := loc.Meters(0.2)
		for i, a := range anchors {
			ranges[i] = state.AnchorReport{Range: j.Distance(a) + loc.Meters(r.Float64())*err - err/2}
		}

		p, acc := ep.ComputePosition(ranges)

		if i%100 == 0 {
			fmt.Println(i)
		}
		err = loc.Meters(math.Max(float64(err), 0.01))
		err *= 2

		fmt.Printf("p %v, accuracy +/-%.2f, actual dist %f, \n", p, float64(acc*3)*math.Sqrt(3), p.Distance(j))
		if p.Distance(j) > err {
			t.Fatalf("idx %d, j %v, p %v, accuracy rms %f, accuracy +/- %.2f, actual dist %f", i, j, p, acc, float64(acc*3)*math.Sqrt(3), p.Distance(j))
		}
	}
}
