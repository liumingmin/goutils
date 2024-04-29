package container

import "sync"

type Pool[T any] struct {
	internal *sync.Pool
}

func NewPool[T any](f func() T) Pool[T] {
	return Pool[T]{
		internal: &sync.Pool{
			New: func() interface{} {
				return f()
			},
		},
	}
}

func (p Pool[T]) Get() T {
	return p.internal.Get().(T)
}

func (p Pool[T]) Put(t T) {
	p.internal.Put(t)
}
