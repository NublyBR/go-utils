package resolver

import "errors"

var (
	ErrTimeout = errors.New("operation timed out")
	ErrClosed  = errors.New("resolver closed")
)

type Resolver[T any] interface {
	Read() (T, error)
	Write(T) error
	Close() error
}

type ResolverBytes interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
