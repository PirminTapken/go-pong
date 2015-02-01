package main

import (
	"runtime"
)

// Thread is a single OS thread
type Thread struct {
	exec   chan func() interface{}
	result <-chan interface{}
	quit   chan chan error
}

// Thread returns the handle to a new OS thread
func NewThread() *Thread {
	t := &Thread{}
	t.exec = make(chan func() interface{})
	// have result local to be able to send
	r := make(chan interface{})
	t.result = r
	t.quit = make(chan chan error)
	go func() {
		runtime.LockOSThread()
		for {
			select {
			case q := <-t.quit:
				runtime.UnlockOSThread()
				q <- nil
				return
			case f := <-t.exec:
				r <- f()
			}
		}
	}()
	return t
}

// Exec executes f in it's thread, returning the result
func (t *Thread) Exec(f func() interface{}) interface{} {
	t.exec <- f
	return <-t.result
}

func (t *Thread) Close() error {
	r := make(chan error)
	t.quit <- r
	return <-r
}
