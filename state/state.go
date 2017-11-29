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
	Power  bool
	Motors [2]Motor
}
