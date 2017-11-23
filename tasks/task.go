package tasks

import "github.com/BoozeBoys/jfino-app/state"

type Task interface {
	Perform(*state.State) error
}
