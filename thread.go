package main

import (
	"runtime"
)

// Thread is a single OS thread
type Thread struct {
	exec   chan func() interface{}
	result <-chan interface{}
	Quit   chan bool
}

// Thread returns the handle to a new OS thread
func NewThread() *Thread {
	t := &Thread{}
	t.exec = make(chan func() interface{})
	// have result local to be able to send
	r := make(chan interface{})
	t.result = r
	t.Quit = make(chan bool)
	go func() {
		runtime.LockOSThread()
		for {
			select {
			case <-t.Quit:
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
