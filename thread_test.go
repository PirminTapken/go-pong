package main

import "testing"
import "time"

func TestThreadExec(t *testing.T) {
	thread := NewThread()
	r := make(chan bool)
	go func() {
		if thread.Exec(func() interface{} {
			return true
		}) != true {
			t.Error("Thread failed to return bool")
		}
		r <- true
	}()
	select {
	case <-r:
		// we were successful
	case <-time.After(1 * time.Millisecond):
		t.Error("Thread didn't respond")
	}
}
