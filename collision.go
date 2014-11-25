package main

import (
	"math"
)

type (
	Point  [2]float64
	Vector [2]float64
)

func (v *Vector) Length() float64 {
	return math.Sqrt(
		v[0]*v[0] + v[1]*v[1],
	)
}

type Object struct {
	Coordinates Point
	Direction   Vector
	Length      float64
	Height      float64
}

func TestIntersect(IntervalA, IntervalB [2]float64) bool {
	return IntervalB[0] <= IntervalA[1] && IntervalA[0] <= IntervalB[1]
}

func VectorAdd(a, b [2]float64) [2]float64 {
	return [2]float64{
		a[0] + b[0],
		a[1] + b[1],
	}
}

func AreLinear(a, b [2]float64) bool {
	var r [2]float64
	r[0] = a[0] / b[0]
	r[1] = a[1] / b[1]
	return r[0] == r[1]
}

// Test if A, B would hit in next Interval defined by
// Speed of respective object
// Return reflection angle
func HitTest(A, B Object) (float64, float64) {
	AP := [2]float64{
		A.Coordinates[0],
		A.Coordinates[1],
	}
	AD := [2]float64{
		A.Direction[0],
		A.Direction[1],
	}
	BP := [2]float64{
		B.Coordinates[0],
		B.Coordinates[1],
	}
	BD := [2]float64{
		B.Direction[0],
		B.Direction[1],
	}
	if AreLinear(AD, BD) {
		// same direction
		// doesn't mean they collide
		return 0, 0
	}
	if AP == BP {
		return 0, 0
	}
	return 0, 0
}
