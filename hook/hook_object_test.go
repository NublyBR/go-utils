package hook

import (
	"io"
	"testing"
)

func TestHookObject(t *testing.T) {
	type Target struct{}
	type Event struct{}

	var (
		hook = NewObject[*Target]()

		target = Target{}
		event  = Event{}

		called bool
		calls  int
		err    error
	)

	hook.Add(func(t *Target, e Event) {
		called = true
	})

	calls, err = hook.Run(&target, event)
	if !called {
		t.Error("Expected hook to be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHookObjectError(t *testing.T) {
	type Target struct{}
	type Event struct{}

	var (
		hook = NewObject[Target]()
	)

	hook.Add(func(_ Target, _ Event) error {
		return io.EOF
	})

	_, err := hook.Run(Target{}, Event{})
	if err != io.EOF {
		t.Errorf("Expected error %v, got %v", io.EOF, err)
	}
}
