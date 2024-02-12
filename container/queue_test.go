package container

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEnqueueBack(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueBack(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueBack dequeue front: ", deItem)
		}
	}
	t.Log("queue is empty ", q.IsEmpty())
	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

}

func TestDequeueFront(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueBack(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueBack dequeue front: ", deItem)
		}
	}

	for i := 0; i < 5; i++ {
		deItem := q.DequeueFront()
		if deItem != "" {
			t.Log("DequeueFront dequeue front: ", deItem)
		}
	}

	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 5; i++ {
		deItem := q.DequeueFront()
		if deItem != "" {
			t.Log("DequeueFront dequeue front: ", deItem)
		}
	}

	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

	t.Log("queue is empty ", q.IsEmpty())
}

func TestEnqueueFront(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueFront(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueFront dequeue back: ", deItem)
		}
	}
	t.Log("queue is empty ", q.IsEmpty())
	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})
}

func TestDequeueBack(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueFront(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueFront dequeue back: ", deItem)
		}
	}

	for i := 0; i < 5; i++ {
		deItem := q.DequeueBack()
		if deItem != "" {
			t.Log("DequeueBack dequeue back: ", deItem)
		}
	}
	t.Log("queue is empty ", q.IsEmpty())

	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

	for i := 0; i < 5; i++ {
		deItem := q.DequeueBack()
		if deItem != "" {
			t.Log("DequeueBack dequeue back: ", deItem)
		}
	}

	t.Log("queue is empty ", q.IsEmpty())
}

func TestQueueClear(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueBack(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueBack dequeue front: ", deItem)
		}
	}
	t.Log("queue is empty ", q.IsEmpty())
	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

	q.Clear()
	t.Log("queue is empty ", q.IsEmpty())
	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})
}

func TestQueueFindBy(t *testing.T) {
	q := InitQueue()
	t.Log("queue is empty ", q.IsEmpty())

	for i := 0; i < 25; i++ {
		deItem := q.EnqueueBack(fmt.Sprint(i))
		if deItem != "" {
			t.Log("EnqueueBack dequeue front: ", deItem)
		}
	}

	q.Range(func(i string) bool {
		t.Log("left:", i)
		return true
	})

	i20 := q.FindOneBy(func(i string) bool {
		ni, _ := strconv.Atoi(i)
		if ni == 20 {
			return true
		}
		return false
	})
	if i20 != "20" {
		t.FailNow()
	} else {
		t.Log("FindOneBy: ", i20)
	}

	data := q.FindBy(func(i string) bool {
		ni, _ := strconv.Atoi(i)
		if ni%3 == 0 {
			return true
		}
		return false
	})
	t.Log("modern 3 :", data)
}

func InitQueue() *Queue[string] {
	return NewQueue[string](10)
}
