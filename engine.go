package main

import ()

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDLE_LENGTH    = 0.2
	PADDLE_WIDTH     = 0.05
	PADDLE_COLOR     = uint32(0xffffff)
)

// Engine is our little nice graphics engine
type Engine struct {
	// Thread is public so it can be used by other
	// goroutines that need to do stuff in the sdl
	// thread
	Thread *Thread
}

// Close closes the engine
// error is always nil and just there to match
// io.Closer
func (e *Engine) Close() error {
	err := e.Thread.Close()
	return err
}

// NewEngine creates the engine.
// This basically creates the background texture and stores it away
func NewEngine(windowName string, X, Y, W, H int) (e *Engine, err error) {
	e = &Engine{Thread: NewThread()}
	return e, nil
}
