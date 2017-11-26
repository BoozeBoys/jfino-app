package tasks

import (
	"math"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
)

type EstimatePosition struct {
	anchors map[uint]loc.Point // anchors location
}

func NewEstimatePosition(anchors map[uint]loc.Point) *EstimatePosition {
	return &EstimatePosition{anchors: anchors}
}

func (ep *EstimatePosition) Perform(s *state.State) error {

	return nil
}

/*ErrorRms gives the Root Mean Square of the differences between
 * the measured anchor ranges and the distance from the point j to each anchor.
 * TODO: use anchor power report to weight the mean computation.
 */
func (ep *EstimatePosition) ErrorRms(report map[uint]state.AnchorReport, j loc.Point) float64 {
	e := 0.0

	for id, r := range report {
		dist := j.Distance(ep.anchors[id])
		e += math.Pow(dist-r.Range, 2)
	}

	return math.Sqrt(e / float64(len(report)))
}

/*FindBoundingBox finds the box that contains the target point we are looking for.
 * Given the anchors position and range reports, compute the maximum box where
 * we need to search for the robot position.
 */
func (ep *EstimatePosition) FindBoundingBox(report map[uint]state.AnchorReport) loc.Box {
	//find working space
	id := uint(0)
	//find a valid anchor id
	for k := range report {
		id = k
		break
	}
	minX := ep.anchors[id].X
	minY := ep.anchors[id].Y
	minZ := ep.anchors[id].Z
	maxX, maxY, maxZ := minX, minY, minZ

	for id, r := range report {
		x := ep.anchors[id].X
		y := ep.anchors[id].Y
		z := ep.anchors[id].Z

		findMin := func(min, v float64) float64 {
			m := v - r.Range
			if m < min {
				return m
			}
			return min
		}

		findMax := func(max, v float64) float64 {
			m := v + r.Range
			if m > max {
				return m
			}
			return max
		}

		minX = findMin(minX, x)
		minY = findMin(minY, y)
		minZ = findMin(minZ, z)

		maxX = findMax(maxX, x)
		maxY = findMax(maxY, y)
		maxZ = findMax(maxZ, z)
	}

	p0 := loc.Point{X: minX, Y: minY, Z: minZ}
	p1 := loc.Point{X: maxX, Y: maxY, Z: maxZ}
	return loc.Box{P0: p0, P1: p1}
}

func (ep *EstimatePosition) ComputePosition(report map[uint]state.AnchorReport) (j loc.Point, best float64) {
	bbox := ep.FindBoundingBox(report)
	center := bbox.Center()

	return center, 0
}
