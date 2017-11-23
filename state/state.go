package state

type Motor struct {
	Speed   int
	Current int
}

type State struct {
	Power  bool
	Motors [2]Motor
}
