package resolver

import "sync"

type resolverPool[I comparable, T any] struct {
	mp map[I]Resolver[T]
	mu sync.Mutex
}

func NewPool[I comparable, T any]() ResolverPool[I, T] {
	return &resolverPool[I, T]{
		mp: make(map[I]Resolver[T]),
	}
}

func (r *resolverPool[I, T]) Clean() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k := range r.mp {
		if !r.mp[k].Alive() {
			delete(r.mp, k)
		}
	}
}

func (r *resolverPool[I, T]) Put(i I, rs Resolver[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.mp[i] = rs
}

func (r *resolverPool[I, T]) Read(i I) (T, error) {
	r.mu.Lock()
	v, ok := r.mp[i]
	r.mu.Unlock()

	if ok {
		return v.Read()
	}

	var zero T
	return zero, ErrNotFound
}

func (r *resolverPool[I, T]) Write(i I, value T) error {
	r.mu.Lock()
	v, ok := r.mp[i]
	if ok {
		if _, isSingle := v.(*resolverSingle[T]); isSingle {
			delete(r.mp, i)
		}
	}
	r.mu.Unlock()

	if ok {
		return v.Write(value)
	}

	return ErrNotFound
}

func (r *resolverPool[I, T]) Close(i I) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.mp[i]; ok {
		delete(r.mp, i)
		return v.Close()
	}

	return ErrNotFound
}

func (r *resolverPool[I, T]) Alive(i I) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.mp[i]; ok {
		return v.Alive()
	}

	return false
}
