package main

import (
	"testing"
)

func getNewSDLEngine(t *testing.T, s string, a, b, c, d int) (e *SDLEngine) {
	e, err := NewSDLEngine(s, a, b, c, d)
	if err != nil {
		t.Error(err)
	}
	if e == nil {
		t.Error("Engine is nil")
	}
	return
}

func TestEngineClose(t *testing.T) {
	e := getNewSDLEngine(t, "test", 100, 100, 100, 100)
	err := e.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestEngineSetTitle(t *testing.T) {
	s := "TestTitle"
	e := getNewSDLEngine(t, "test", 100, 100, 100, 100)
	e.SetTitle(s)
	if e.Title() != s {
		t.Error("setting title failed")
	}
	err := e.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestAppendChild(t *testing.T) {
	r := new(Node)
	n := []*Node{new(Node), new(Node), new(Node)}
	for _, c := range n {
		r.AppendChild(c)
	}
	if n[0] != r.FirstChild ||
		n[1] != n[0].NextSibling ||
		n[1] != n[2].PrevSibling ||
		n[2] != r.LastChild ||
		n[1].Parent != r {
		t.Error("Fail")
	}
}

func TestRemoveChild(t *testing.T) {
	r := new(Node)
	n := []*Node{new(Node), new(Node), new(Node)}
	for _, c := range n {
		r.AppendChild(c)
	}
	r.RemoveChild(n[1])
	if n[0].NextSibling != n[2] ||
		n[2].PrevSibling != n[0] {
		t.Error("Fail")
	}
}
