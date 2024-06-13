package hook

import (
	"io"
	"reflect"
	"testing"
)

func TestHook(t *testing.T) {
	type A struct{}
	type B struct{}

	var (
		hook = New()

		aCalled, bCalled bool

		calls int
		err   error
	)

	hook.Add(func(_ A) {
		aCalled = true
	}, func(_ B) {
		bCalled = true
	})

	calls, err = hook.Run(A{})
	if !aCalled || bCalled {
		t.Error("Expected A hook to be called and B hook to not be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	aCalled, bCalled = false, false

	calls, err = hook.Run(B{})
	if !bCalled || aCalled {
		t.Error("Expected B hook to be called and A hook to not be called")
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHookMulti(t *testing.T) {
	type A struct{}

	var (
		hook = New()

		num, calls int
		err        error
	)

	for i := 0; i < 10; i++ {
		hook.Add(func(_ A) {
			num++
		})
	}

	calls, err = hook.Run(A{})
	if num != 10 {
		t.Errorf("Expected 10 calls, got %d", num)
	}

	if calls != 10 {
		t.Errorf("Expected 10 calls, got %d", calls)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHookError(t *testing.T) {
	type A struct{}

	var (
		hook = New()
	)

	hook.Add(func(_ A) error {
		return io.EOF
	})

	_, err := hook.Run(A{})
	if err != io.EOF {
		t.Errorf("Expected error %v, got %v", io.EOF, err)
	}
}

func TestHookHandlers(t *testing.T) {
	var h = New()

	h.Add(func(int) {}, func(string) {}, func(float32) {})

	var (
		handlers = h.Handlers()
		expect   = h.(*hook).mp
	)

	if reflect.ValueOf(handlers).UnsafePointer() == reflect.ValueOf(expect).UnsafePointer() {
		t.Error("Expected handlers to be copied")
	}

	if !reflect.DeepEqual(handlers, expect) {
		t.Errorf("Expected %v, got %v", expect, handlers)
	}
}
