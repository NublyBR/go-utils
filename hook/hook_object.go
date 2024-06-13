package hook

import (
	"reflect"
	"sync"
)

type HookObject[T any] interface {
	// Add adds a function to the hook
	//
	//   hook.Add(func(target T, e Event) { ... })
	Add(fn any)

	// Run calls all the functions in the hook
	//
	//   hook.Run(target, Event{})
	Run(target T, i any) int
}

type hookObject[T any] struct {
	mp map[reflect.Type][]reflect.Value
	mu sync.Mutex
}

func NewObject[T any]() HookObject[T] {
	return &hookObject[T]{
		mp: make(map[reflect.Type][]reflect.Value),
	}
}

func (h *hookObject[T]) Add(fn any) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("expected a function")
	}

	if t.NumIn() != 2 || t.NumOut() != 0 {
		panic("expected a function with two inputs and no output")
	}

	if t.In(0) != reflect.TypeOf((*T)(nil)).Elem() {
		panic("expected the first input to be the target type")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.mp[t.In(1)] = append(h.mp[t.In(1)], reflect.ValueOf(fn))
}

func (h *hookObject[T]) Run(target T, i any) int {
	h.mu.Lock()
	v := h.mp[reflect.TypeOf(i)]
	h.mu.Unlock()

	for _, f := range v {
		f.Call([]reflect.Value{reflect.ValueOf(target), reflect.ValueOf(i)})
	}

	return len(v)
}
