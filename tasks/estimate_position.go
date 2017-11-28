package tasks

import (
	"math"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
)

type EstimatePosition struct {
	anchors map[int]loc.Point // anchors location
}

func NewEstimatePosition(anchors map[int]loc.Point) *EstimatePosition {
	return &EstimatePosition{anchors: anchors}
}

func (ep *EstimatePosition) Perform(s *state.State) error {

	return nil
}

/*ErrorRms gives the Root Mean Square of the differences between
 * the measured anchor ranges and the distance from the point j to each anchor.
 * TODO: use anchor power report to weight the mean computation.
 */
func (ep *EstimatePosition) ErrorRms(report map[int]state.AnchorReport, j loc.Point) float64 {
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
func (ep *EstimatePosition) FindBoundingBox(report map[int]state.AnchorReport) loc.Box {
	//find working space
	id := 0
	//find a valid anchor id
	for k := range report {
		id = k
		break
	}

	//initialize
	var pmin loc.Point
	var pmax loc.Point
	for i, v := range ep.anchors[id] {
		pmin[i] = v
		pmax[i] = pmin[i]
	}

	//find max points
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

	return loc.Box{pmin, pmax}
}

func (ep *EstimatePosition) ComputePosition(report map[int]state.AnchorReport) (j loc.Point, accuracy float64) {
	box := ep.FindBoundingBox(report)
	accuracy = box[0].Distance(box[1])
	maxIter := 50
loop:
	for i := 0; i < maxIter; i++ {
		s := box.Bisect()
		accuracy = box[0].Distance(box[1])

		for _, v := range s {
			if acc := ep.ErrorRms(report, v.Center()); acc < accuracy {
				accuracy = acc
				box = v.Expand(1.525)
				if accuracy <= 0.005*math.Sqrt(3)/3 { // +/-0.5 cm
					break loop
				}
			}
		}
	}

	return box.Center(), accuracy
}
