package loc_test

import "testing"
import "github.com/BoozeBoys/jfino-app/loc"
import "math"

func TestPointIsEqual(t *testing.T) {
	p0 := loc.Point{0, 0, 0}
	p1 := p0
	p2 := loc.Point{1, 1, 1}
	p3 := loc.Point{0123e+12, 1231.30, 4324324.453454350}

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
	p4[0] -= 1e-14

	if !p3.IsEqual(p4) {
		t.FailNow()
	}

	p4[0] -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}

	p4[0] += 1e-2
	p4[1] -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}

	p4[1] += 1e-2
	p4[2] -= 1e-2
	if p3.IsEqual(p4) {
		t.FailNow()
	}
}

func TestPointDistance(t *testing.T) {
	p1 := loc.Point{0, 0, 0}
	p2 := loc.Point{1, 1, 1}

	if math.Abs(p1.Distance(p2)-math.Sqrt(3)) > 1e-6 {
		t.FailNow()
	}

	if math.Abs(p2.Distance(p1)-math.Sqrt(3)) > 1e-6 {
		t.FailNow()
	}

	p3 := loc.Point{-1, -1, -1}

	if math.Abs(p2.Distance(p1)-p1.Distance(p3)) > 1e-6 {
		t.FailNow()
	}

	p4 := loc.Point{-1, 1, -1}

	if math.Abs(p2.Distance(p4)-math.Sqrt(8)) > 1e-6 {
		t.FailNow()
	}
}

func TestBoxCenter(t *testing.T) {
	p0 := loc.Point{0, 0, 0}
	p1 := loc.Point{112, 1e23, 1e-12}
	b := loc.Box{P0: p0, P1: p1}

	if pc := b.Center(); !pc.IsEqual(loc.Point{56, 5e22, 5e-13}) {
		t.FailNow()
	}
}

func TestBoxTranspose(t *testing.T) {
	b := loc.Box{P0: loc.Point{0, 0, 0}, P1: loc.Point{4, 4, 4}}

	bx := b.TransposeX()
	if !bx.P0.IsEqual(loc.Point{4, 0, 0}) {
		t.FailNow()
	}
	if !bx.P1.IsEqual(loc.Point{0, 4, 4}) {
		t.FailNow()
	}

	by := b.TransposeY()
	if !by.P0.IsEqual(loc.Point{0, 4, 0}) {
		t.FailNow()
	}
	if !by.P1.IsEqual(loc.Point{4, 0, 4}) {
		t.FailNow()
	}

	bz := b.TransposeZ()
	if !bz.P0.IsEqual(loc.Point{0, 0, 4}) {
		t.FailNow()
	}
	if !bz.P1.IsEqual(loc.Point{4, 4, 0}) {
		t.FailNow()
	}
}

func TestBoxisEqual(t *testing.T) {
	b := loc.Box{P0: loc.Point{0, 0, 0}, P1: loc.Point{4, 4, 4}}
	b1 := loc.Box{P0: loc.Point{0, 0, 0}, P1: loc.Point{4, 4, 4}}

	if !b.IsEqual(b1) {
		t.Fatalf("b %v, b1 %v", b, b1)
	}

	b1.P0[0] = 4
	b1.P1[0] = 0
	if !b.IsEqual(b1) {
		t.Fatalf("b %v, b1 %v", b, b1)
	}

	b1.P0[1] = 4
	b1.P1[1] = 0
	if !b.IsEqual(b1) {
		t.Fatalf("b %v, b1 %v", b, b1)
	}

	b1.P0[2] = 4
	b1.P1[2] = 0
	if !b.IsEqual(b1) {
		t.Fatalf("b %v, b1 %v", b, b1)
	}

	b1.P0[0] = 2
	b1.P1[0] = 2
	if b.IsEqual(b1) {
		t.Fatalf("b %v, b1 %v", b, b1)
	}

}

func TestBoxBisect(t *testing.T) {
	p0 := loc.Point{0, 0, 0}
	p1 := loc.Point{8, 8, 8}
	b := loc.Box{P0: p0, P1: p1}

	bcheck := []loc.Box{
		loc.Box{P0: loc.Point{0, 0, 0}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{0, 8, 0}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{8, 8, 0}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{8, 0, 0}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{0, 0, 8}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{0, 8, 8}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{8, 8, 8}, P1: loc.Point{4, 4, 4}},
		loc.Box{P0: loc.Point{8, 0, 8}, P1: loc.Point{4, 4, 4}},
	}

	s := b.Bisect()
	for i, v := range bcheck {
		if !s[i].IsEqual(v) {
			t.Fatalf("s %v, check %v", s[i], v)
		}
	}
}
