package cache

import (
	"testing"
	"time"
)

func TestCacheExpire(t *testing.T) {
	t.Parallel()

	var (
		c = New[int, int]()
	)

	c.Put(5*time.Millisecond, 0, 1337)

	if v, ok := c.Get(0); !ok {
		t.Error("Expected true, got false")
	} else if v != 1337 {
		t.Errorf("Expected 1337, got %d", v)
	}

	time.Sleep(5 * time.Millisecond)

	if _, ok := c.Get(0); ok {
		t.Error("Expected false, got true")
	}
}

func TestCacheRemove(t *testing.T) {
	t.Parallel()

	var (
		c = New[int, int]()
	)

	c.Put(5*time.Millisecond, 0, 1337)

	if v, ok := c.Remove(0); !ok {
		t.Error("Expected true, got false")
	} else if v != 1337 {
		t.Errorf("Expected 1337, got %d", v)
	}

	if _, ok := c.Remove(0); ok {
		t.Error("Expected false, got true")
	}
}
