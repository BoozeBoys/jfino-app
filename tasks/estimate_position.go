package tasks

import (
	"log"
	"math"

	"github.com/BoozeBoys/jfino-app/loc"

	"github.com/BoozeBoys/jfino-app/state"
	"gonum.org/v1/gonum/optimize"
)

type AnchorCfg struct {
	Loc    loc.Point
	Offset loc.Meters
}

type EstimatePosition struct {
	anchors map[string]AnchorCfg // anchors location
}

func NewEstimatePosition(anchors map[string]AnchorCfg) *EstimatePosition {
	return &EstimatePosition{anchors: anchors}
}

/*Perform computes the current position of the robot.
 * TODO: filter the position, discard wrong positions,
 * use a smaller box when we have a fix.
 */
func (ep *EstimatePosition) Perform(s *state.State) error {
	box := ep.FindBoundingBox(s.RangeReport)

	/* Look for position around the floor (+20cm from the floor) */
	box[0][2] = 0
	box[1][2] = 0.2
	p, acc := ep.ComputePosition(box, s.RangeReport)
	s.CurrentPosition = p
	s.PositionAccuracy = acc
	return nil
}

func (ep *EstimatePosition) ErrorRms(report map[string]state.AnchorReport, j loc.Point) loc.Meters {
	mse := 0.0

	for id, r := range report {
		dist := float64(j.Distance(ep.anchors[id].Loc))
		mse += math.Pow(dist-float64(r.Range-ep.anchors[id].Offset), 2)
	}

	return loc.Meters(mse / float64(len(report)))
}

/*ErrorPtoP gives the Peak-to-Peak (+/-) error between
 * the measured anchor ranges and the distance from the point j to each anchor.
 * TODO: use anchor power report to weight the mean computation.
 */
func (ep *EstimatePosition) ErrorPtoP(report map[string]state.AnchorReport, j loc.Point) loc.Meters {

	// return error as distance +/- from average (stddev *3 which covers 99.7% probability)
	return loc.Meters(math.Sqrt(float64(ep.ErrorRms(report, j))) * 3)
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
	for i, v := range ep.anchors[id].Loc {
		pmin[i] = v
		pmax[i] = pmin[i]
	}

	//find max points
	for id, r := range report {
		for i, v := range ep.anchors[id].Loc {
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
	f := func(x []float64) float64 {
		j = loc.Point{loc.Meters(x[0]), loc.Meters(x[1]), loc.Meters(x[2])}
		return float64(ep.ErrorRms(report, j))
	}

	grad := func(grad []float64, x []float64) {
		j = loc.Point{loc.Meters(x[0]), loc.Meters(x[1]), loc.Meters(x[2])}

		for i := range j {
			grad[i] = 0
		}

		for id, r := range report {
			dr := float64(r.Range)
			d := float64(j.Distance(ep.anchors[id].Loc))
			for i := range j {
				xi := float64(ep.anchors[id].Loc[i])
				dv := (x[i] - xi) / d
				grad[i] += (2 * dv * (d - dr))
			}
		}

		for i := range j {
			grad[i] /= float64(len(report))
		}
	}

	p := optimize.Problem{
		Func: f,
		Grad: grad,
	}
	mid := box.Center()
	x := []float64{float64(mid[0]), float64(mid[1]), float64(mid[2])}
	settings := optimize.DefaultSettingsLocal()
	settings.Recorder = nil
	settings.GradientThreshold = 1e-12
	settings.FunctionConverge = nil

	result, err := optimize.Minimize(p, x, settings, &optimize.BFGS{})
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}

	return loc.Point{loc.Meters(result.X[0]), loc.Meters(result.X[1]), loc.Meters(result.X[2])}, loc.Meters(result.F)
}
