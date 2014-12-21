package main

import (
	"fmt"
	"math"
)

// no inplace operations
type Vector2d [2]float64

// Add vector w to v and return result
func (v *Vector2d) Add(w *Vector2d) *Vector2d {
	u := new(Vector2d)
	for i, _ := range v {
		u[i] = v[i] + w[i]
	}
	return u
}

// Scale vector by factor f
func (v *Vector2d) Scale(f float64) *Vector2d {
	u := new(Vector2d)
	for i, _ := range v {
		u[i] = v[i] * f
	}
	return u
}

func (v *Vector2d) Dot(w *Vector2d) float64 {
	return v[0]*w[0] + v[1]*w[1]
}

func (v *Vector2d) Sub(w *Vector2d) *Vector2d {
	return v.Add(w.Scale(-1.0))
}

func (v *Vector2d) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

// Copy returns a pointer to a new copy of v
func (v *Vector2d) Copy() *Vector2d {
	return &Vector2d{v[0], v[1]}
}

// Reflect reflects w from v. Assumption is made that they
// intersect. This is not checked on purpose
func (v *Vector2d) Reflect(w *Vector2d) *Vector2d {
	n := &Vector2d{-v[1], v[0]}
	r := w.Sub(n.Scale((2 * w.Dot(n))))
	return r
}

// Inverse returns the inverse vector in regard to the
// dot product
func (v *Vector2d) Inverse() *Vector2d {
	return &Vector2d{-v[1], v[0]}
}

func (v *Vector2d) GoString() string {
	return fmt.Sprintf(`&Vector2d{%v, %v}`, v[0], v[1])
}

func (v *Vector2d) String() string {
	return fmt.Sprintf(`[%v, %v]`, v[0], v[1])
}

type Line [2]*Vector2d

// get the factor for intersection with line k
func (l *Line) Intersect(k *Line) float64 {
	// see http://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
	// for the calculation
	A, B := l[0], l[1]
	C, D := k[0], k[1]

	E := B.Sub(A)
	F := D.Sub(C)

	P := E.Inverse()
	h := (A.Sub(C)).Dot(P) / F.Dot(P)

	return h
}

func (l *Line) Vector2d() *Vector2d {
	return l[1].Sub(l[0])
}

func (l *Line) GoString() string {
	return fmt.Sprintf(`&Line{%v, %v}`, l[0], l[1])
}

func (l *Line) String() string {
	return fmt.Sprintf("Line(%v, %v)", l[0], l[1])
}
