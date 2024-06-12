package resolver

import (
	"sync"
	"testing"
	"time"
)

func testResolverSingle(timeout, waitWriter, waitReader time.Duration, valWrite int) (valRead int, errRead error, errWrite error) {
	var (
		res = NewSingle[int](timeout)
		wg  sync.WaitGroup
	)

	go func() {
		time.Sleep(waitWriter)
		errWrite = res.Write(valWrite)
		wg.Done()
	}()
	wg.Add(1)

	go func() {
		time.Sleep(waitReader)
		valRead, errRead = res.Read()
		wg.Done()
	}()
	wg.Add(1)

	wg.Wait()

	return
}

func testExpect(t *testing.T, valRead, valReadExpect int, errRead, errReadExpect error, errWrite, errWriteExpect error) {
	if valRead != valReadExpect {
		t.Errorf("Expected value from Read: %d, got %d", valReadExpect, valRead)
	}

	if errRead != errReadExpect {
		t.Errorf("Expected error from Read: %v, got %v", errReadExpect, errRead)
	}

	if errWrite != errWriteExpect {
		t.Errorf("Expected error from Write: %v, got %v", errWriteExpect, errWrite)
	}
}

const (
	tHigh = 9 * time.Millisecond
	tMed  = 6 * time.Millisecond
	tLow  = 3 * time.Millisecond
)

func TestSingleReadFirst(t *testing.T) {
	t.Parallel()

	valRead, errRead, errWrite := testResolverSingle(tHigh, tMed, tLow, 42)

	testExpect(t,
		valRead, 42,
		errRead, nil,
		errWrite, nil,
	)
}

func TestSingleWriteFirst(t *testing.T) {
	t.Parallel()

	valRead, errRead, errWrite := testResolverSingle(tHigh, tLow, tMed, 42)

	testExpect(t,
		valRead, 42,
		errRead, nil,
		errWrite, nil,
	)
}

func TestSingleTimeoutWrite(t *testing.T) {
	t.Parallel()

	valRead, errRead, errWrite := testResolverSingle(tMed, tHigh, tLow, 42)

	testExpect(t,
		valRead, 0,
		errRead, ErrTimeout,
		errWrite, ErrTimeout,
	)
}

func TestSingleTimeoutRead(t *testing.T) {
	t.Parallel()

	valRead, errRead, errWrite := testResolverSingle(tMed, tLow, tHigh, 42)

	testExpect(t,
		valRead, 0,
		errRead, ErrTimeout,
		errWrite, ErrTimeout,
	)
}
