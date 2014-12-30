package main

import (
	"testing"
	"time"
)

func TestWallIntersections(t *testing.T) {
	walls := []*Line{
		&Line{&Vector2d{0, 0}, &Vector2d{0, 1}},
		&Line{&Vector2d{0, 1}, &Vector2d{1, 1}},
		&Line{&Vector2d{1, 1}, &Vector2d{1, 0}},
		&Line{&Vector2d{1, 0}, &Vector2d{0, 0}},
	}
	line := &Line{
		&Vector2d{0.1, 0.5},
		&Vector2d{0.4, 0.6},
	}
	intersections := WallIntersections(walls, line)
	for _, intersection := range intersections {
		if intersection.IntersectAt > 0 && intersection.IntersectAt < 1 {
			t.Errorf("intersection we have, captain")
		}
	}
	line = &Line{
		&Vector2d{0.5, 0.5},
		&Vector2d{0.5, 1.5},
	}
	intersections = WallIntersections(walls, line)
	for i, intersection := range intersections {
		if i == 2 {
			if !(intersection.IntersectAt > 0) ||
				!(intersection.IntersectAt < 1) {
				t.Errorf(
					"%v should intersect %v",
					line,
					intersection.Wall,
				)
			}
		} else {
			if intersection.IntersectAt > 0 &&
				intersection.IntersectAt < 1 {
				t.Errorf("intersection we have, captain")
			}
		}
	}
}

func TestUpdateBall(t *testing.T) {
	fakeUniverse := map[string]Object{
		"Ball": Object{
			W:  0.1,
			H:  0.1,
			X:  0.5,
			Y:  0.5,
			DX: 0.1,
			DY: 0.1,
		},
	}
	fakeUniverseChan := make(chan map[string]Object, 1)
	fakeUniverseChan <- fakeUniverse
	errChan := make(chan error, 1)
	d := time.Second
	UpdateBall(fakeUniverseChan, errChan, d)
	select {
	case err := <-errChan:
		t.Error("UpdateBall caused error", err)
	default:
		// continue
	}
	fakeUniverse = <-fakeUniverseChan
	pos := &Vector2d{fakeUniverse["Ball"].X, fakeUniverse["Ball"].Y}
	should := &Vector2d{0.6, 0.6}
	if !pos.Equals(should) {
		t.Errorf(
			"Ball(%v, %v) is not at [%v, %v]",
			fakeUniverse["Ball"].X,
			fakeUniverse["Ball"].Y,
			0.6, 0.6,
		)
	}
}

func TestUpdateBallWall(t *testing.T) {
	fakeUniverse := map[string]Object{
		"Ball": Object{
			X:  0.5,
			Y:  0.9,
			DX: 0.1,
			DY: 0.2,
		},
	}
	fakeUniverseChan := make(chan map[string]Object, 1)
	fakeUniverseChan <- fakeUniverse
	errChan := make(chan error, 1)
	d := time.Second
	UpdateBall(fakeUniverseChan, errChan, d)
	select {
	case err := <-errChan:
		t.Error("UpdateBall caused error", err)
	default:
		// continue
	}
	fakeUniverse = <-fakeUniverseChan
	if !(fakeUniverse["Ball"].X == 0.6 && fakeUniverse["Ball"].Y == 0.9) {
		t.Errorf(
			"Ball(%v, %v) is not at [%v, %v]",
			fakeUniverse["Ball"].X,
			fakeUniverse["Ball"].Y,
			0.6, 0.6,
		)
	}

}
