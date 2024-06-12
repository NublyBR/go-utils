package resolver

import "sync"

type resolverPoolBytes[I comparable] struct {
	mp map[I]ResolverBytes
	mu sync.Mutex
}

func NewPoolBytes[I comparable]() ResolverPoolBytes[I] {
	return &resolverPoolBytes[I]{
		mp: make(map[I]ResolverBytes),
	}
}

func (r *resolverPoolBytes[I]) Clean() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k := range r.mp {
		if !r.mp[k].Alive() {
			delete(r.mp, k)
		}
	}
}

func (r *resolverPoolBytes[I]) Put(i I, rs ResolverBytes) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.mp[i] = rs
}

func (r *resolverPoolBytes[I]) Read(i I, b []byte) (int, error) {
	r.mu.Lock()
	v, ok := r.mp[i]
	r.mu.Unlock()

	if ok {
		return v.Read(b)
	}

	return 0, ErrNotFound
}

func (r *resolverPoolBytes[I]) Write(i I, b []byte) (int, error) {
	r.mu.Lock()
	v, ok := r.mp[i]
	r.mu.Unlock()

	if ok {
		return v.Write(b)
	}

	return 0, ErrNotFound
}

func (r *resolverPoolBytes[I]) Close(i I) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.mp[i]; ok {
		delete(r.mp, i)
		return v.Close()
	}

	return ErrNotFound
}

func (r *resolverPoolBytes[I]) Alive(i I) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v, ok := r.mp[i]; ok {
		return v.Alive()
	}

	return false
}
