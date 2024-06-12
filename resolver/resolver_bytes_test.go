package resolver

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestOrder(t *testing.T) {
	t.Parallel()

	var (
		expect = []byte(strings.Repeat("Hello, World!", 1024))

		input  = bytes.NewReader(expect)
		output = bytes.NewBuffer(nil)

		b  = NewBytes(time.Second)
		wg sync.WaitGroup
	)

	go func() {
		var buf = make([]byte, 5)

		n, err := io.CopyBuffer(output, b, buf)
		if err != nil {
			t.Error(err)
		}

		if int(n) != len(expect) {
			t.Errorf("Expected %d bytes, got %d", len(expect), n)
		}

		wg.Done()
	}()
	wg.Add(1)

	go func() {
		var buf = make([]byte, 4)

		n, err := io.CopyBuffer(b, input, buf)
		if err != nil {
			t.Error(err)
		}

		if int(n) != len(expect) {
			t.Errorf("Expected %d bytes, got %d", len(expect), n)
		}

		if err := b.Close(); err != nil {
			t.Error(err)
		}

		wg.Done()
	}()
	wg.Add(1)

	wg.Wait()

	if !bytes.Equal(output.Bytes(), expect) {
		t.Errorf("Expected %q, got %q", expect, output.Bytes())
	}
}
