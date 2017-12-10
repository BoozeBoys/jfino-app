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
	Power int        // dBm
}

type State struct {
	ConnectedToFW    bool
	Power            bool
	Motors           [2]Motor
	RangeReport      map[int]AnchorReport
	CurrentPosition  loc.Point
	PositionAccuracy loc.Meters
}
