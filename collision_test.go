package main

import (
	"math"
	"testing"
)

// Test collision between two objects,
// one moves straight down against the other,
// which itself is at rest
func TestCollisionDown(t *testing.T) {
	A := Object{
		Coordinates: Point{5, 5},
		Direction:   Vector{1, 0},
		Length:      1,
		Height:      1,
	}
	B := Object{
		Coordinates: Point{7, 7},
		Direction:   Vector{0, 0},
		Length:      2,
		Height:      2,
	}

	ReflectionA, ReflectionB := HitTest(A, B)
	if ReflectionA != 0 {
		t.Error("A is not reflected properly")
	}
	if ReflectionB != 0 {
		t.Error("B is not reflected properly")
	}
}

func TestAreLinear(t *testing.T) {
	a := [2]float64{1.0, 2.0}
	b := [2]float64{2.0, 4.0}
	if !AreLinear(a, b) {
		t.Error("a and b should be linear!")
	}
	b = [2]float64{3.0, 1.0}
	if AreLinear(a, b) {
		t.Error("a and b should *not* be linear!")
	}
}

func TestCollisionDiagonal(t *testing.T) {
	A := Object{
		Coordinates: Point{2, 5},
		Direction:   Vector{1, 1},
		Length:      1,
		Height:      1,
	}
	B := Object{
		Coordinates: Point{8, 10},
		Direction:   Vector{0, 0},
		Length:      2,
		Height:      2,
	}

	ReflectionA, ReflectionB := HitTest(A, B)
	if ReflectionA != math.Pi/2 {
		t.Error("A isnt reflected at 90 degrees, but instead rad:", ReflectionA)
	}
	if ReflectionB != 0 {
		t.Error("B shouldn't be reflected at all!")
	}
}
