package safego

import (
	"context"

	"github.com/liumingmin/goutils/log"
)

type Handler func(err interface{})

var DefaultHandler = func(err interface{}) {
	log.ErrorStack(context.Background(), "goroutine: ", err)
}

func Go(f func(), handler ...Handler) {
	handle := DefaultHandler
	switch len(handler) {
	case 1:
		handle = handler[0]
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				handle(r)
			}
		}()

		f()
	}()
}

func GoWithArgs(f func(args ...interface{}), args ...interface{}) {
	GoWithHandler(f, DefaultHandler, args...)
}

// GoWithHandler will run function f(args... interface{}) in a go routines.
// And handler the panic if occur with given handler.
func GoWithHandler(f func(args ...interface{}), handler Handler, args ...interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				handler(r)
			}
		}()

		f(args...)
	}()
}
