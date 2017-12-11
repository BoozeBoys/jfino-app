package state

import (
	"github.com/BoozeBoys/jfino-app/loc"
)

type Motor struct {
	Speed   int
	Current int
}

type AnchorReport struct {
	Range loc.Meters // m
	Power float64    // dBm
}

type State struct {
	Power            bool
	Motors           [2]Motor
	RangeReport      map[string]AnchorReport
	CurrentPosition  loc.Point
	PositionAccuracy loc.Meters
}
