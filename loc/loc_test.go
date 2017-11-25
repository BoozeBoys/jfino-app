package loc_test

import "testing"
import "github.com/BoozeBoys/jfino-app/loc"
import "math"

func TestIsEqual(t *testing.T) {
	p0 := loc.Point{X: 0, Y: 0, Z: 0}
	p1 := p0
	p2 := loc.Point{X: 1, Y: 1, Z: 1}
	p3 := loc.Point{X: 0123e+12, Y: 1231.30, Z: 4324324.453454350}

	if p1.IsEqual(p2) {
		t.FailNow()
	}

	p1 = p2
	if !p1.IsEqual(p2) {
		t.FailNow()
	}

	if p0.IsEqual(p1) {
		t.FailNow()
	}

	p4 := p3
	if !p3.IsEqual(p4) {
		t.FailNow()
	}
	p4.X -= 1e-14

	if !p3.IsEqual(p4) {
		t.FailNow()
	}

	p4.X -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}

	p4.X += 1e-2
	p4.Y -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}

	p4.Y += 1e-2
	p4.Z -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}
}

func TestDistance(t *testing.T) {
	p1 := loc.Point{X: 0, Y: 0, Z: 0}
	p2 := loc.Point{X: 1, Y: 1, Z: 1}

	if math.Abs(p1.Distance(p2)-math.Sqrt(3)) > 1e-6 {
		t.FailNow()
	}

	if math.Abs(p2.Distance(p1)-math.Sqrt(3)) > 1e-6 {
		t.FailNow()
	}

	p3 := loc.Point{X: -1, Y: -1, Z: -1}

	if math.Abs(p2.Distance(p1)-p1.Distance(p3)) > 1e-6 {
		t.FailNow()
	}

	p4 := loc.Point{X: -1, Y: 1, Z: -1}

	if math.Abs(p2.Distance(p4)-math.Sqrt(8)) > 1e-6 {
		t.FailNow()
	}
}

func TestBoxCenter(t *testing.T) {
	p0 := loc.Point{X: 0, Y: 0, Z: 0}
	p1 := loc.Point{X: 112, Y: 1e23, Z: 1e-12}
	b := loc.Box{P0: p0, P1: p1}

	if pc := b.Center(); !pc.IsEqual(loc.Point{X: 56, Y: 5e22, Z: 5e-13}) {
		t.FailNow()
	}
}

func TestBoxSlice(t *testing.T) {
	p0 := loc.Point{X: 0, Y: 0, Z: 0}
	p1 := loc.Point{X: 8, Y: 8, Z: 8}
	b := loc.Box{P0: p0, P1: p1}

	pcheck := loc.Box{P0: loc.Point{X: 0, Y: 0, Z: 0}, P1: loc.Point{X: 0, Y: 8, Z: 0}}

	s := b.Slice()
	if !s[0].P0.IsEqual(pcheck.P0) {
		t.FailNow()
	}

	//TODO
}
