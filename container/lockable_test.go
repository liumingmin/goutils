package container

import (
	"sync"
	"testing"
)

func TestLockable(t *testing.T) {
	var i Lockable[int]
	i.Set(100)
	if i.Get() != 100 {
		t.Error(i.Get())
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		i.Update(func(i int) int { return i + 1 })
		wg.Done()
	}()
	go func() {
		i.Update(func(i int) int { return i + 2 })
		wg.Done()
	}()
	go func() {
		i.Update(func(i int) int { return i - 3 })
		wg.Done()
	}()
	wg.Wait()

	if i.Get() != 100 {
		t.Error(i.Get())
	}
}
