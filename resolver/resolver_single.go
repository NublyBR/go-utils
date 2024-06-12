package resolver

import (
	"sync"
	"time"
)

type resolverSingle[T any] struct {
	ticker *time.Ticker
	chData chan T

	isClosed bool
	chClose  chan struct{}

	mu sync.Mutex
}

// NewSingle creates a new Resolver with the given timeout.
func NewSingle[T any](timeout time.Duration) Resolver[T] {
	return &resolverSingle[T]{
		ticker: time.NewTicker(timeout),

		chData: make(chan T),

		chClose: make(chan struct{}),
	}
}

// Read reads the value or returns a timeout error if it cannot read within the given timeout.
func (r *resolverSingle[T]) Read() (T, error) {
	select {
	case <-r.chClose:
		var zero T
		return zero, ErrClosed

	case value := <-r.chData:
		return value, nil

	case <-r.ticker.C:
		var zero T
		return zero, ErrTimeout
	}
}

// Write writes the value or returns a timeout error if it cannot write within the given timeout.
func (r *resolverSingle[T]) Write(value T) error {
	select {
	case <-r.chClose:
		return ErrClosed

	case r.chData <- value:
		return nil

	case <-r.ticker.C:
		return ErrTimeout
	}
}

// Close closes the resolver.
func (r *resolverSingle[T]) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isClosed {
		return ErrClosed
	}
	r.isClosed = true

	close(r.chClose)
	r.ticker.Stop()
	return nil
}
