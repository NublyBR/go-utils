package resolver

import (
	"testing"
	"time"
)

func TestMultiWriteTimeout(t *testing.T) {
	t.Parallel()

	var multi = NewMulti[int](time.Millisecond)

	if err := multi.Write(42); err != ErrTimeout {
		t.Errorf("Expected ErrTimeout, got %v", err)
	}
}

func TestMultiReadTimeout(t *testing.T) {
	t.Parallel()

	var multi = NewMulti[int](time.Millisecond)

	if _, err := multi.Read(); err != ErrTimeout {
		t.Errorf("Expected ErrTimeout, got %v", err)
	}
}
