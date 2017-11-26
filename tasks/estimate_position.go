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

	//initialize
	var pmin loc.Point
	var pmax loc.Point
	for i := 0; i < len(ep.anchors[id]); i++ {
		pmin[i] = ep.anchors[id][i]
		pmax[i] = pmin[i]
	}

	for id, r := range report {
		for i, v := range ep.anchors[id] {
			min := v - r.Range
			max := v + r.Range
			if min < pmin[i] {
				pmin[i] = min
			}

			if max > pmax[i] {
				pmax[i] = max
			}
		}
	}

	return loc.Box{P0: pmin, P1: pmax}
}

func (ep *EstimatePosition) ComputePosition(report map[uint]state.AnchorReport) (j loc.Point, accuracy float64) {
	box := ep.FindBoundingBox(report)
	accuracy = box.P0.Distance(box.P1)
	maxIter := 25
	id := 0

loop:
	for i := 0; i < maxIter; i++ {
		s := box.Bisect()

		for i, v := range s {
			if acc := ep.ErrorRms(report, v.Center()); acc < accuracy {
				accuracy = acc
				id = i
			}

			if accuracy <= 0.01 { // 1 cm
				break loop
			}
		}
		box = s[id]
	}

	return box.Center(), accuracy
}
