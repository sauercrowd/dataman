package definitions

import (
	"time"
)

type DataReceiver interface {
	AddTuples(tp []ProviderTuple)
	AddTuple(tp ProviderTuple)
	Finish()
}

type ProviderTuple struct {
	Date    time.Time
	ID      string
	Author  string
	Content string
	Replies int64
	Score   int64
}
