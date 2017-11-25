package loc

import (
	"math"
)

type Point struct {
	X, Y, Z float64 // meters
}

/*IsEqual returns true if two points are the same */
func (p *Point) IsEqual(p1 Point) bool {
	eps := 1e-64
	if math.Abs(p.X-p1.X) > eps {
		return false
	}

	if math.Abs(p.Y-p1.Y) > eps {
		return false
	}

	if math.Abs(p.Z-p1.Z) > eps {
		return false
	}

	return true
}

/*Distance computes the distance between two points */
func (p *Point) Distance(p1 Point) float64 {
	x := p1.X - p.X
	y := p1.Y - p.Y
	z := p1.Z - p.Z

	return math.Sqrt(x*x + y*y + z*z)
}
