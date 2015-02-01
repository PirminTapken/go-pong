package main

import (
	"testing"
)

func TestThreadExec(t *testing.T) {
	thread := NewThread()
	if thread.Exec(func() interface{} {
		return true
	}) != true {
		thread.Close()
		t.Error("Thread failed to return bool")
	}
}
