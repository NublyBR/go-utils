package resolver

import (
	"sync"
	"time"
)

type resolverMulti[T any] struct {
	timeout time.Duration
	chData  chan T

	isClosed bool
	chClose  chan struct{}

	mu sync.Mutex
}

// NewSingle creates a new Resolver with the given timeout.
func NewMulti[T any](timeout time.Duration) Resolver[T] {
	return &resolverMulti[T]{
		timeout: timeout,

		chData: make(chan T),
	}
}

// Read reads the value or returns a timeout error if it cannot read within the given timeout.
func (r *resolverMulti[T]) Read() (T, error) {
	var ticker = time.NewTicker(r.timeout)
	defer ticker.Stop()

	select {
	case <-r.chClose:
		var zero T
		return zero, ErrClosed

	case value := <-r.chData:
		return value, nil

	case <-ticker.C:
		var zero T
		return zero, ErrTimeout
	}
}

// Write writes the value or returns a timeout error if it cannot write within the given timeout.
func (r *resolverMulti[T]) Write(value T) error {
	var ticker = time.NewTicker(r.timeout)
	defer ticker.Stop()

	select {
	case <-r.chClose:
		return ErrClosed

	case r.chData <- value:
		return nil

	case <-ticker.C:
		return ErrTimeout
	}
}

// Close closes the resolver.
func (r *resolverMulti[T]) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isClosed {
		return ErrClosed
	}
	r.isClosed = true

	close(r.chClose)
	return nil
}
