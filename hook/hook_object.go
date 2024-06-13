package hook

import (
	"reflect"
	"sync"
)

type HookObject[T any] interface {
	// Add adds a function to the hook
	//
	//   hook.Add(func(target T, e Event) { ... })
	Add(fn ...any)

	// Run calls all the functions in the hook
	//
	//   hook.Run(target, Event{})
	Run(target T, obj any) (int, error)
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

func (h *hookObject[T]) add(fn any) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("expected a function")
	}

	if t.NumIn() != 2 || t.NumOut() > 1 {
		panic("expected a function with two inputs and an optional error output")
	}

	if t.In(0) != reflect.TypeOf((*T)(nil)).Elem() {
		panic("expected the first input to be the target type")
	}

	if t.NumOut() == 1 && t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("expected an error output")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.mp[t.In(1)] = append(h.mp[t.In(1)], reflect.ValueOf(fn))
}

func (h *hookObject[T]) Add(fn ...any) {
	for _, f := range fn {
		h.add(f)
	}
}

func (h *hookObject[T]) Run(target T, obj any) (int, error) {
	h.mu.Lock()
	v := h.mp[reflect.TypeOf(obj)]
	h.mu.Unlock()

	for i, f := range v {
		var res = f.Call([]reflect.Value{reflect.ValueOf(target), reflect.ValueOf(obj)})

		if len(res) > 0 {
			return i + 1, res[0].Interface().(error)
		}
	}

	return len(v), nil
}
