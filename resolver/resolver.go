package resolver

import "errors"

var (
	ErrTimeout  = errors.New("operation timed out")
	ErrClosed   = errors.New("resolver closed")
	ErrNotFound = errors.New("resolver not found")
)

type Resolver[T any] interface {
	// Read reads the value or returns a timeout error if it cannot read within the given timeout.
	Read() (T, error)
	// Write writes the value or returns a timeout error if it cannot write within the given timeout.
	Write(T) error
	// Close closes the resolver.
	Close() error
	// Alive checks if the resolver is still alive.
	Alive() bool
}

type ResolverBytes interface {
	// Read reads the value or returns a timeout error if it cannot read within the given timeout.
	Read([]byte) (int, error)
	// Write writes the value or returns a timeout error if it cannot write within the given timeout.
	Write([]byte) (int, error)
	// Close closes the resolver.
	Close() error
	// Alive checks if the resolver is still alive.
	Alive() bool
}

type ResolverPool[I comparable, T any] interface {
	// Read calls the Read method of the given resolver.
	Read(I) (T, error)
	// Write calls the Write method of the given resolver.
	Write(I, T) error
	// Close closes the given resolver.
	Close(I) error
	// Alive checks if the given resolver is still alive.
	Alive(I) bool
	// Put adds the given resolver to the pool with the given key.
	Put(I, Resolver[T])
	// Clean cleans the pool of dead resolvers.
	Clean()
}

type ResolverPoolBytes[I comparable] interface {
	// Read calls the Read method of the given bytes resolver.
	Read(I, []byte) (int, error)
	// Write calls the Write method of the given bytes resolver.
	Write(I, []byte) (int, error)
	// Close closes the given bytes resolver.
	Close(I) error
	// Alive checks if the given bytes resolver is still alive.
	Alive(I) bool
	// Put adds the given bytes resolver to the pool with the given key.
	Put(I, ResolverBytes)
	// Clean cleans the pool of dead byte resolvers.
	Clean()
}
