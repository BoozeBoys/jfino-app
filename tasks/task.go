package tasks

type Motor struct {
	Speed   int
	Current int
}

type State struct {
	Power  bool
	Motors [2]Motor
}

type Task interface {
	Perform(*State) error
}
