package hook

import "testing"

func TestHookObject(t *testing.T) {
	type Target struct{}
	type Event struct{}

	var (
		hook = NewObject[*Target]()

		target = Target{}
		event  = Event{}

		called bool
		calls  int
	)

	hook.Add(func(t *Target, e Event) {
		called = true
	})

	calls = hook.Run(&target, event)
	if !called {
		t.Error("Expected hook to be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}
