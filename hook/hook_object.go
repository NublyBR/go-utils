package hook

import (
	"reflect"
	"sync"
)

type HookObject[T any] interface {
	// Add adds a function to the hook, with optional target input and error output
	//
	//   hook.Add(func(target T, e Event) error { ... })
	//   hook.Add(func(target T, e Event) { ... })
	//   hook.Add(func(e Event) error { ... })
	//   hook.Add(func(e Event) { ... })
	Add(fn ...any)

	// Run calls all the functions in the hook
	//
	//   hook.Run(target, Event{})
	Run(target T, obj any) (int, error)

	// Returns a copy of the handlers map
	Handlers() map[reflect.Type][]reflect.Value
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

	if t.NumIn() == 0 || t.NumIn() > 2 || t.NumOut() > 1 {
		panic("expected a function with one to two inputs and an optional error output")
	}

	var et reflect.Type

	if t.NumIn() == 1 {
		et = t.In(0)
	} else {
		et = t.In(1)
		if t.In(0) != reflect.TypeOf((*T)(nil)).Elem() {
			panic("expected the first input to be the target type")
		}
	}

	if t.NumOut() == 1 && t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("expected an error output")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.mp[et] = append(h.mp[et], reflect.ValueOf(fn))
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
		var res []reflect.Value

		if f.Type().NumIn() == 1 {
			res = f.Call([]reflect.Value{reflect.ValueOf(obj)})
		} else {
			res = f.Call([]reflect.Value{reflect.ValueOf(target), reflect.ValueOf(obj)})
		}

		if len(res) > 0 {
			return i + 1, res[0].Interface().(error)
		}
	}

	return len(v), nil
}

func (h *hookObject[T]) Handlers() map[reflect.Type][]reflect.Value {
	var m = make(map[reflect.Type][]reflect.Value, len(h.mp))

	h.mu.Lock()
	defer h.mu.Unlock()

	for k, v := range h.mp {
		var vc = make([]reflect.Value, len(v))
		copy(vc, v)
		m[k] = vc
	}

	return m
}
