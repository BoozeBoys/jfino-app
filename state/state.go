package state

type Motor struct {
	Speed   int
	Current int
}

type AnchorReport struct {
	Range float64 // m
	Power int     // dBm
}

type State struct {
	Power  bool
	Motors [2]Motor
}
