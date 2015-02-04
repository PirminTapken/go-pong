package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDLE_LENGTH    = 0.2
	PADDLE_WIDTH     = 0.05
	PADDLE_COLOR     = uint32(0xffffff)
)

// Engine is our little nice graphics engine
type Engine struct {
	cleanupFns []func()
	window     *sdl.Window
	renderer   *sdl.Renderer
	// Thread is public so it can be used by other
	// goroutines that need to do stuff in the sdl
	// thread
	Thread *Thread
}

// NewEngine creates the engine.
// This basically creates the background texture and stores it away
func NewEngine(windowName string, X, Y, W, H int) (e *Engine, err error) {
	e = &Engine{Thread: NewThread(), cleanupFns: make([]func(), 0)}
	err = e.sdlInit()
	if err != nil {
		return e, err
	}
	err = e.CreateWindowAndRenderer(W, H, 0)
	e.SetTitle(windowName)
	return e, err
}

// Close closes the engine
// error is always nil and just there to match
// io.Closer
func (e *Engine) Close() error {
	e.cleanup()
	err := e.Thread.Close()
	return err
}

func (e *Engine) Title() string {
	return e.Thread.Exec(func() interface{} {
		return e.window.GetTitle()
	}).(string)
}

func (e *Engine) SetTitle(s string) {
	_ = e.Thread.Exec(func() interface{} {
		e.window.SetTitle(s)
		return nil
	})
}

// init calls sdl init in sdl thread
func (e *Engine) sdlInit() error {
	r := e.Thread.Exec(func() interface{} {
		// This is neccessary otherwise nil error can't be converted to
		// interface and back somehow...
		e := struct{ err error }{}
		if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
			e.err = sdl.GetError()
		}
		return e
	}).(struct{ err error })
	e.cleanupFns = append(e.cleanupFns, sdl.Quit)
	return r.err
}

// cleanup cleans everything up
func (e *Engine) cleanup() {
	for i := len(e.cleanupFns); i > 0; i-- {
		e.Thread.Exec(func() interface{} {
			e.cleanupFns[i-1]()
			return nil
		})
	}
}

func (e *Engine) CreateWindowAndRenderer(w, h int, flags uint32) error {
	type resp struct {
		w *sdl.Window
		r *sdl.Renderer
		e error
	}
	r := e.Thread.Exec(func() interface{} {
		w, r, e := sdl.CreateWindowAndRenderer(w, h, flags)
		return resp{w: w, r: r, e: e}
	}).(resp)
	if r.e != nil {
		return r.e
	}
	e.window = r.w
	e.renderer = r.r
	e.cleanupFns = append(e.cleanupFns, e.window.Destroy)
	e.cleanupFns = append(e.cleanupFns, e.renderer.Destroy)
	return nil
}
