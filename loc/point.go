package loc

import (
	"math"
)

type Meters float64

type Point [3]Meters // x, y, z coords in meters

/*IsEqual returns true if two points are the same */
func (p *Point) IsEqual(p1 Point) bool {
	eps := 1e-64
	for i := range p {
		if math.Abs(float64(p[i]-p1[i])) > eps {
			return false
		}
	}

	return true
}

/*Distance computes the distance between two points */
func (p *Point) Distance(p1 Point) Meters {
	sum := 0.0
	for i := range p {
		sum += math.Pow(float64(p1[i]-p[i]), 2)
	}

	return Meters(math.Sqrt(sum))
}
