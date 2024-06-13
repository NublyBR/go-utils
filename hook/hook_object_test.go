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

		num, calls int
		err        error
	)

	hook.Add(func(t *Target, e Event) {
		num++
	}, func(e Event) {
		num++
	})

	calls, err = hook.Run(&target, event)
	if num != 2 {
		t.Errorf("Expected both functions to be called, got %d", num)
	}

	if calls != 2 {
		t.Errorf("Expected 2 calls, got %d", calls)
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
