package tasks

import (
	"math"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
)

type EstimatePosition struct {
	anchors map[string]loc.Point // anchors location
}

func NewEstimatePosition(anchors map[string]loc.Point) *EstimatePosition {
	return &EstimatePosition{anchors: anchors}
}

/*Perform computes the current position of the robot.
 * TODO: filter the position, discard wrong positions,
 * use a smaller box when we have a fix.
 */
func (ep *EstimatePosition) Perform(s *state.State) error {
	box := ep.FindBoundingBox(s.RangeReport)
	p, acc := ep.ComputePosition(box, s.RangeReport)
	s.CurrentPosition = p
	s.PositionAccuracy = acc
	return nil
}

/*ErrorPtoP gives the Peak-to-Peak (+/-) error between
 * the measured anchor ranges and the distance from the point j to each anchor.
 * TODO: use anchor power report to weight the mean computation.
 */
func (ep *EstimatePosition) ErrorPtoP(report map[string]state.AnchorReport, j loc.Point) loc.Meters {
	e := 0.0

	for id, r := range report {
		dist := float64(j.Distance(ep.anchors[id]))
		e += math.Pow(dist-float64(r.Range), 2)
	}

	// return error as distance +/- from average (stddev *3 which covers 99.7% probability)
	return loc.Meters(math.Sqrt(e/float64(len(report))) * 3)
}

/*FindBoundingBox finds the box that contains the target point we are looking for.
 * Given the anchors position and range reports, compute the maximum box where
 * we need to search for the robot position.
 */
func (ep *EstimatePosition) FindBoundingBox(report map[string]state.AnchorReport) loc.Box {
	//find working space
	id := ""
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

const minAccuracy = loc.Meters(0.005) // +/-0.5 cm
const expandFactor = 1.55
const maxIter = 50

func (ep *EstimatePosition) ComputePosition(box loc.Box, report map[string]state.AnchorReport) (j loc.Point, accuracy loc.Meters) {

loop:
	for i := 0; i < maxIter; i++ {
		s := box.Bisect()
		accuracy = box[0].Distance(box[1])

		for _, v := range s {
			if acc := ep.ErrorPtoP(report, v.Center()); acc < accuracy {
				accuracy = acc
				box = v.Expand(expandFactor)
				if accuracy <= minAccuracy {
					break loop
				}
			}
		}
	}

	return box.Center(), accuracy
}
