package main

import (
	"testing"
)

func TestEngineClose(t *testing.T) {
	e, err := NewEngine("test", 100, 100, 100, 100)
	if err != nil {
		t.Error(err)
	}
	if e == nil {
		t.Error("Engine is nil")
	}
	err = e.Close()
	if err != nil {
		t.Error(err)
	}
}
