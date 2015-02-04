package main

import (
	"testing"
)

func getNewEngine(t *testing.T, s string, a, b, c, d int) (e *Engine) {
	e, err := NewEngine(s, a, b, c, d)
	if err != nil {
		t.Error(err)
	}
	if e == nil {
		t.Error("Engine is nil")
	}
	return
}

func TestEngineClose(t *testing.T) {
	e := getNewEngine(t, "test", 100, 100, 100, 100)
	err := e.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestEngineSetTitle(t *testing.T) {
	s := "TestTitle"
	e := getNewEngine(t, "test", 100, 100, 100, 100)
	e.SetTitle(s)
	if e.Title() != s {
		t.Error("setting title failed")
	}
	err := e.Close()
	if err != nil {
		t.Error(err)
	}
}
