package hook

import "testing"

func TestHook(t *testing.T) {
	type A struct{}
	type B struct{}

	var (
		hook = New()

		aCalled, bCalled bool
		calls            int
	)

	hook.Add(func(_ A) {
		aCalled = true
	})

	hook.Add(func(_ B) {
		bCalled = true
	})

	calls = hook.Run(A{})
	if !aCalled || bCalled {
		t.Error("Expected A hook to be called and B hook to not be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}

	aCalled, bCalled = false, false

	calls = hook.Run(B{})
	if !bCalled || aCalled {
		t.Error("Expected B hook to be called and A hook to not be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}

func TestHookMulti(t *testing.T) {
	type A struct{}

	var (
		hook = New()

		num, calls int
	)

	for i := 0; i < 10; i++ {
		hook.Add(func(_ A) {
			num++
		})
	}

	calls = hook.Run(A{})
	if num != 10 {
		t.Errorf("Expected 10 calls, got %d", num)
	}

	if calls != 10 {
		t.Errorf("Expected 10 calls, got %d", calls)
	}
}
