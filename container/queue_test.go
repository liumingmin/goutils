package container

import (
	"fmt"
	"testing"
)

func TestEnqueueBack(t *testing.T) {
	q := InitQueue()
	if !q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueBack(fmt.Sprint(i))

		if i >= q.Len() {
			expectedItem := fmt.Sprint(i - q.Len())
			if deItem != expectedItem {
				t.Error(deItem)
			}
		} else {
			if deItem != "" {
				t.Error(deItem)
			}
		}
	}

	if q.IsEmpty() {
		t.FailNow()
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })

}

func TestDequeueFront(t *testing.T) {
	q := InitQueue()
	if !q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 25; i++ {
		q.EnqueueBack(fmt.Sprint(i))
	}

	for i := 0; i < 5; i++ {
		deItem := q.DequeueFront()
		expectedItem := fmt.Sprint(i + 25 - 10)
		if deItem != expectedItem {
			t.Error(deItem)
		}
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })

	if q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		deItem := q.DequeueFront()
		expectedItem := fmt.Sprint(i + 25 - 10 + 5)
		if deItem != expectedItem {
			t.Error(deItem)
		}
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })

	if !q.IsEmpty() {
		t.FailNow()
	}
}

func TestEnqueueFront(t *testing.T) {
	q := InitQueue()
	if !q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueFront(fmt.Sprint(i))
		if i >= q.Len() {
			expectedItem := fmt.Sprint(i - q.Len())
			if deItem != expectedItem {
				t.Error(deItem)
			}
		} else {
			if deItem != "" {
				t.Error(deItem)
			}
		}
	}

	if q.IsEmpty() {
		t.FailNow()
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })
}

func TestDequeueBack(t *testing.T) {
	q := InitQueue()
	if !q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 25; i++ {
		q.EnqueueFront(fmt.Sprint(i))
	}

	for i := 0; i < 5; i++ {
		deItem := q.DequeueBack()
		expectedItem := fmt.Sprint(i + 25 - 10)
		if deItem != expectedItem {
			t.Error(deItem)
		}
	}
	if q.IsEmpty() {
		t.FailNow()
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })

	for i := 0; i < 5; i++ {
		deItem := q.DequeueBack()
		expectedItem := fmt.Sprint(i + 25 - 10 + 5)
		if deItem != expectedItem {
			t.Error(deItem)
		}
	}

	if !q.IsEmpty() {
		t.FailNow()
	}
}

func TestQueueClear(t *testing.T) {
	q := InitQueue()
	if q.Cap() != 10 {
		t.FailNow()
	}

	if !q.IsEmpty() {
		t.FailNow()
	}

	for i := 0; i < 25; i++ {
		q.EnqueueBack(fmt.Sprint(i))
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })

	q.Clear()

	if q.Cap() != 10 {
		t.FailNow()
	}

	if !q.IsEmpty() {
		t.FailNow()
	}

	// q.Range(func(i string) bool {
	// 	t.Log("left:", i)
	// 	return true
	// })
}

func TestQueueFindBy(t *testing.T) {
	q := NewQueue[int](25)

	for i := 0; i < 25; i++ {
		q.EnqueueBack(i)
	}

	i20 := q.FindOneBy(func(i int) bool {
		return i == 20
	})

	if i20 != 20 {
		t.FailNow()
	}

	items := q.FindBy(func(i int) bool {
		return i%3 == 0
	})

	for _, item := range items {
		if item%3 != 0 {
			t.Error(item)
		}
	}
}

func TestQueueRange(t *testing.T) {
	q := NewQueue[int](10)
	for i := 0; i < 25; i++ {
		q.EnqueueBack(i)
	}

	j := 15
	q.Range(func(i int) bool {
		if i != j {
			t.Error(i)
		}
		j++
		return true
	})
}

func InitQueue() *Queue[string] {
	return NewQueue[string](10)
}
