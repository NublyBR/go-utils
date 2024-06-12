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

func NewMulti[T any](timeout time.Duration) Resolver[T] {
	return &resolverMulti[T]{
		timeout: timeout,

		chData: make(chan T),

		chClose: make(chan struct{}),
	}
}

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

func (r *resolverMulti[T]) Alive() bool {
	select {
	case <-r.chClose:
		return false
	default:
		return true
	}
}
