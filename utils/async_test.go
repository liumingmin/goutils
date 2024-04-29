package utils

import (
	"context"
	"testing"
	"time"
)

func TestAsyncInvokeWithTimeout(t *testing.T) {
	f1 := false
	f2 := false
	result := AsyncInvokeWithTimeout(time.Second*1, func() {
		time.Sleep(time.Millisecond * 500)
		f1 = true
	}, func() {
		time.Sleep(time.Millisecond * 500)
		f2 = true
	})

	if !result {
		t.FailNow()
	}

	if !f1 {
		t.FailNow()
	}

	if !f2 {
		t.FailNow()
	}
}

func TestAsyncInvokeWithTimeouted(t *testing.T) {
	f1 := false
	f2 := false
	result := AsyncInvokeWithTimeout(time.Second*1, func() {
		time.Sleep(time.Millisecond * 1500)
		f1 = true
	}, func() {
		time.Sleep(time.Millisecond * 500)
		f2 = true
	})

	if result {
		t.FailNow()
	}

	if f1 {
		t.FailNow()
	}

	if !f2 {
		t.FailNow()
	}
}

func TestAsyncInvokesWithTimeout(t *testing.T) {
	f1 := false
	f2 := false

	fns := []func(){
		func() {
			time.Sleep(time.Millisecond * 500)
			f1 = true
		}, func() {
			time.Sleep(time.Millisecond * 500)
			f2 = true
		},
	}
	result := AsyncInvokesWithTimeout(time.Second*1, fns)

	if !result {
		t.FailNow()
	}

	if !f1 {
		t.FailNow()
	}

	if !f2 {
		t.FailNow()
	}
}

func TestSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t1 := time.Now()
	go func() {
		cancel()
	}()

	Sleep(ctx, 2*time.Second)

	if time.Since(t1) > time.Second {
		t.FailNow()
	}
}
