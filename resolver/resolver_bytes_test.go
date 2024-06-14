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

func TestEmpty(t *testing.T) {
	t.Parallel()

	var b = NewBytes(time.Second)

	go func() {
		var seq = [][]byte{
			nil,
			[]byte("Hello"),
			nil,
			[]byte(", "),
			nil,
			[]byte("World!"),
			nil,
		}
		for _, data := range seq {
			var n, err = b.Write(data)
			if err != nil {
				t.Error(err)
			}

			if int(n) != len(data) {
				t.Errorf("Expected %d bytes, got %d", len(data), n)
			}
		}

		if err := b.Close(); err != nil {
			t.Error(err)
		}
	}()

	var buf = bytes.NewBuffer(nil)

	var _, err = io.Copy(buf, b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf.Bytes(), []byte("Hello, World!")) {
		t.Errorf("Expected %q, got %q", []byte("Hello, World!"), buf.Bytes())
	}
}
