package hook

import (
	"reflect"
	"sync"
)

type Hook interface {
	// Add adds a function to the hook
	//
	//   hook.Add(func(e Event) { ... })
	Add(fn any)

	// Run calls all the functions in the hook
	//
	//   hook.Run(Event{})
	Run(i any) int
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

func (h *hook) Add(fn any) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("expected a function")
	}

	if t.NumIn() != 1 || t.NumOut() != 0 {
		panic("expected a function with one input and no output")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.mp[t.In(0)] = append(h.mp[t.In(0)], reflect.ValueOf(fn))
}

func (h *hook) Run(i any) int {
	h.mu.Lock()
	v := h.mp[reflect.TypeOf(i)]
	h.mu.Unlock()

	for _, f := range v {
		f.Call([]reflect.Value{reflect.ValueOf(i)})
	}

	return len(v)
}
