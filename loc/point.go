package loc

import (
	"math"
)

type Point [3]float64 // x, y, z coords in meters

/*IsEqual returns true if two points are the same */
func (p *Point) IsEqual(p1 Point) bool {
	eps := 1e-64
	for i := range p {
		if math.Abs(p[i]-p1[i]) > eps {
			return false
		}
	}

	return true
}

/*Distance computes the distance between two points */
func (p *Point) Distance(p1 Point) float64 {
	sum := 0.0
	for i := range p {
		sum += math.Pow(p1[i]-p[i], 2)
	}

	return math.Sqrt(sum)
}
