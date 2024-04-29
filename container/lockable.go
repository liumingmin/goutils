package container

import "sync"

type Lockable[T any] struct {
	sync.Mutex
	value T
}

func (l *Lockable[T]) Get() T {
	l.Lock()
	defer l.Unlock()

	return l.value
}

func (l *Lockable[T]) Set(v T) {
	l.Lock()
	defer l.Unlock()

	l.value = v
}

func (l *Lockable[T]) Update(f func(T) T) {
	l.Lock()
	defer l.Unlock()

	l.value = f(l.value)
}
