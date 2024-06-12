package resolver

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	t.Parallel()

	var (
		pool = NewPool[int, int]()

		orderA = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		orderB = []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

		wg sync.WaitGroup
	)

	pool.Put(0, NewMulti[int](time.Second))
	go func() {
		defer wg.Done()
		defer pool.Close(0)

		for _, i := range orderA {
			if err := pool.Write(0, i); err != nil {
				t.Error(err)
				break
			}
		}
	}()
	wg.Add(1)

	pool.Put(1, NewMulti[int](tHigh))
	go func() {
		defer wg.Done()
		defer pool.Close(1)

		for _, i := range orderB {
			if err := pool.Write(1, i); err != nil {
				t.Error(err)
				break
			}
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()

		var r = make([]int, 0, 10)

		for i := 0; i < 10; i++ {
			v, err := pool.Read(0)
			if err != nil {
				t.Error(err)
				break
			}
			r = append(r, v)
		}

		if !reflect.DeepEqual(orderA, r) {
			t.Errorf("expected %v, got %v", orderA, r)
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()

		var r = make([]int, 0, 10)

		for i := 0; i < 10; i++ {
			v, err := pool.Read(1)
			if err != nil {
				t.Error(err)
				break
			}
			r = append(r, v)
		}

		if !reflect.DeepEqual(orderB, r) {
			t.Errorf("expected %v, got %v", orderA, r)
		}
	}()
	wg.Add(1)

	wg.Wait()
}
