package hook

import (
	"reflect"
	"sync"
)

type Hook interface {
	// Add adds a function to the hook, with optional error output
	//
	//   hook.Add(func(e Event) error { ... })
	//   hook.Add(func(e Event) { ... })
	Add(fn ...any)

	// Run calls all the functions in the hook
	//
	//   hook.Run(Event{})
	Run(obj any) (int, error)
}

type hook struct {
	mp map[reflect.Type][]reflect.Value

	mu sync.Mutex
}

func New() Hook {
	return &hook{
		mp: make(map[reflect.Type][]reflect.Value),
	}
}

func (h *hook) add(fn any) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("expected a function")
	}

	if t.NumIn() != 1 || t.NumOut() > 1 {
		panic("expected a function with one input and an optional error output")
	}

	if t.NumOut() == 1 && t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("expected an error output")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.mp[t.In(0)] = append(h.mp[t.In(0)], reflect.ValueOf(fn))
}

func (h *hook) Add(fn ...any) {
	for _, f := range fn {
		h.add(f)
	}
}

func (h *hook) Run(obj any) (int, error) {
	h.mu.Lock()
	v := h.mp[reflect.TypeOf(obj)]
	h.mu.Unlock()

	for i, f := range v {
		var res = f.Call([]reflect.Value{reflect.ValueOf(obj)})

		if len(res) > 0 {
			return i + 1, res[0].Interface().(error)
		}
	}

	return len(v), nil
}
