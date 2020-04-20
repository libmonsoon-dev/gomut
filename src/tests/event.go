package tests

import "github.com/libmonsoon-dev/jsontest"

// Event is data that produced by go test -json with optional error field
type Event struct {
	jsontest.Event
	Err error
}
