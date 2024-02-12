package container

import (
	"container/list"
	"sync"
)

type Queue[T any] struct {
	items   *list.List
	maxSize int
	mutex   sync.RWMutex
}

func NewQueue[T any](maxSize int) *Queue[T] {
	q := &Queue[T]{
		items:   list.New(),
		maxSize: maxSize,
	}

	q.items.Init()
	return q
}

func (q *Queue[T]) EnqueueBack(item T) (t T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items.PushBack(item)
	if q.items.Len() > q.maxSize {
		return q.dequeueFront()
	}
	return
}

func (q *Queue[T]) DequeueFront() (t T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.dequeueFront()
}

func (q *Queue[T]) dequeueFront() (t T) {
	if q.items.Len() == 0 {
		return
	}

	elem := q.items.Front()
	if elem == nil {
		return
	}

	return q.items.Remove(elem).(T)
}

func (q *Queue[T]) EnqueueFront(item T) (t T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items.PushFront(item)
	if q.items.Len() > q.maxSize {
		return q.dequeueBack()
	}
	return
}

func (q *Queue[T]) DequeueBack() (t T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.dequeueBack()
}

func (q *Queue[T]) dequeueBack() (t T) {
	if q.items.Len() == 0 {
		return
	}

	elem := q.items.Back()
	if elem == nil {
		return
	}

	return q.items.Remove(elem).(T)
}

func (q *Queue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items.Init()
}

func (q *Queue[T]) Range(fn func(t T) bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for elem := q.items.Front(); elem != nil; elem = elem.Next() {
		if !fn(elem.Value.(T)) {
			break
		}
	}
}

func (q *Queue[T]) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.items.Len() == 0
}

func (q *Queue[T]) FindOneBy(fn func(T) bool) (t T) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for elem := q.items.Front(); elem != nil; elem = elem.Next() {
		val := elem.Value.(T)
		if fn(val) {
			return val
		}
	}

	return
}

func (q *Queue[T]) FindBy(fn func(T) bool) (data []T) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for elem := q.items.Front(); elem != nil; elem = elem.Next() {
		val := elem.Value.(T)
		if fn(val) {
			data = append(data, val)
		}
	}

	return
}
