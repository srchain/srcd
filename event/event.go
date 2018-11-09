package event

import (
	"time"
	"sync"
	"reflect"
)

type TypeMuxEvent struct {
	Time time.Time
	Data interface{}
}

// A TypeMux dispatches events to registered receivers. Receivers can be
// registered to handle events of certain type. Any operation
// called after mux is stopped will return ErrMuxClosed.
//
// The zero value is ready to use.
//
// Deprecated: use Feed
type TypeMux struct {
	mutex   sync.RWMutex
	subm    map[reflect.Type][]*TypeMuxSubscription
	stopped bool
}

// TypeMuxSubscription is a subscription established through TypeMux.
type TypeMuxSubscription struct {
	mux     *TypeMux
	created time.Time
	closeMu sync.Mutex
	closing chan struct{}
	closed  bool

	// these two are the same channel. they are stored separately so
	// postC can be set to nil without affecting the return value of
	// Chan.
	postMu sync.RWMutex
	readC  <-chan *TypeMuxEvent
	postC  chan<- *TypeMuxEvent
}

