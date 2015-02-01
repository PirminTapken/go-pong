package main

type Racket struct {
	p      Vector2d
	dim    [2]int // dimension, didn't I have something better?
	update chan func(p Vector2d, dim [2]int)
	get    chan struct {
		p   Vector2d
		dim [2]int
	}
}

func NewRacket() *Racket {
	r := &Racket{}
	r.p = Vector2d{0, 0}
	r.dim = [2]int{0, 0}
	r.update = make(chan func(p Vector2d, dim [2]int))
	r.get = make(chan struct {
		p   Vector2d
		dim [2]int
	})
	return r
}
