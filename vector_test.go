package main

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	a := &Vector2d{3, 3}
	b := &Vector2d{5, 5}
	c := a.Add(b)
	if c[0] != 8 && c[1] != 8 {
		t.Errorf("%v isn't [8, 8]")
	}
}

func TestSub(t *testing.T) {
	a := &Vector2d{3, 3}
	b := &Vector2d{5, 5}
	c := a.Sub(b)
	if c[0] != -2 && c[1] != -2 {
		t.Errorf("%v isn't [-2, -2]", c)
	}
}

func TestLen(t *testing.T) {
	a := &Vector2d{3, 3}
	l := a.Len()
	if l != math.Sqrt(3*3*2) {
		t.Errorf("TestLen failed: %v is not %v", l, 3*3*2)
	}
}

func TestReflect(t *testing.T) {
	a := &Vector2d{3, 3}
	b := &Vector2d{0, 1}
	c := b.Reflect(a)
	if c[0] != -3 || c[1] != 3 {
		t.Errorf("%v is not [3, -3]", c)
	}
}
