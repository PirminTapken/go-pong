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

func TestInverse(t *testing.T) {
	a := &Vector2d{1, 0}
	iA := a.Inverse()
	if a.Dot(iA) != 0 {
		t.Errorf("%v Dot %v isn't 0", a, iA)
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

func TestLineIntersect(t *testing.T) {
	A := &Vector2d{0, 1}
	B := &Vector2d{1, 1}
	AB := &Line{A, B}

	C := &Vector2d{0.5, 0}
	D := &Vector2d{0.5, 2}
	CD := &Line{C, D}

	// h should be:
	// AB == B - A
	// CD == D - C

	h := AB.Intersect(CD)
	if h != 0.5 {
		t.Errorf("h(%v) != 0.5", h)
	}
	k := CD.Intersect(AB)
	if k != 0.5 {
		t.Errorf("h(%v) != 0.5", k)
	}
}

func TestLineVector2d(t *testing.T) {
	A := &Vector2d{0, 1}
	B := &Vector2d{1, 1}
	AB := &Line{A, B}
	vAB := AB.Vector2d()
	if vAB[0] != 1 || vAB[1] != 0 {
		t.Errorf(
			"AB(%v) converted to vAB(%v) isn't [0, 1]",
			AB, vAB,
		)
	}
}
