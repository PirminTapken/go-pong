package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"strconv"
	"time"
)

// CleanupFn is a struct keeping a function with on signature
// and it's name. It's purpose is debugging purposes only
type cleanupFn struct {
	Fn   func()
	Name string
}

// Engine is our little nice graphics engine
type Engine struct {
	renderer      *sdl.Renderer
	Background    *sdl.Texture
	ObjectTexture *sdl.Texture
	window        *sdl.Window
	deferredFns   []cleanupFn
}

// Close closes the engine
// error is always nil and just there to match
// io.Closer
func (e *Engine) Close() error {
	// FIXME cleanup everything
	for i := len(e.deferredFns); i != 0; i-- {
		// TODO maybe handle errors returned
		// by these and aggregate or whatever
		e.deferredFns[i].Fn()
	}
	return nil
}

// StartSdl creates window and renderer and returns a
// cleanup slice that contains functions that should
// be called upon exit to clean up all sdl related stuff
func (e *Engine) startSdl(name string, x, y, w, h int) (
	err error,
) {
	var sdlErr error
	e.deferredFns = make([]cleanupFn, 0, 3)
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		err = fmt.Errorf("sdl.Init failed: %v", sdl.GetError())
		return err
	}
	e.deferredFns = append(e.deferredFns, cleanupFn{sdl.Quit, "sdl.Quit"})
	e.window, sdlErr = sdl.CreateWindow(name, x, y, w, h, 0)
	if sdlErr != nil {
		err = fmt.Errorf("creating sdl.Window failed: %v", err)
		return err
	}
	e.deferredFns = append(e.deferredFns, cleanupFn{e.window.Destroy, "window.Destroy"})
	e.renderer, sdlErr = sdl.CreateRenderer(
		e.window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if sdlErr != nil {
		err = fmt.Errorf("Creating sdl.Renderer failed: %v", err)
		return err
	}
	e.deferredFns = append(e.deferredFns, cleanupFn{e.renderer.Destroy, "renderer.Destroy"})
	return nil
}

// CreateEngine creates the engine.
// This basically creates the background texture and stores it away
func CreateEngine(windowName string, X, Y, W, H int) (e *Engine, err error) {
	e = &Engine{}
	if err = e.startSdl(windowName, X, Y, W, H); err != nil {
		e.Close()
		return e, err
	}
	bgSurface := sdl.CreateRGBSurface(
		0, SCREEN_WIDTH, SCREEN_HEIGHT, 32,
		0, 0, 0, 0,
	)
	if bgSurface == nil {
		e.Close()
		return e, fmt.Errorf("surface creation failed: %v", sdl.GetError())
	}
	if bgSurface.FillRect(nil, BACKGROUND_COLOR) != 0 {
		e.Close()
		return e, fmt.Errorf("filling background failed: %v", sdl.GetError())
	}
	e.Background, err = e.renderer.CreateTextureFromSurface(bgSurface)
	if err != nil {
		return e, fmt.Errorf("Creating BackgroundTexture failed:", err)
	}
	objectSurface := sdl.CreateRGBSurface(0, 1, 1, 32, 0, 0, 0, 0) // just a single color
	objectSurface.FillRect(nil, PADDLE_COLOR)
	e.ObjectTexture, err = e.renderer.CreateTextureFromSurface(objectSurface)
	if err != nil {
		return e, fmt.Errorf("Creating ObjectTexture failed:", err)
	}
	return e, nil
}

// Render the world
func (e *Engine) Render() (err error) {
	e.renderer.Copy(e.Background, nil, nil)
	if err != nil {
		return err
	}
	// TODO what's on with the rest?
	e.renderer.Present()
	return nil
}

// Run the game
// it dirigates all the stuff that needs to be done
func (e *Engine) Run() (err error) {
	fps, err := strconv.Atoi(os.Getenv("FPS"))
	if err != nil {
		fps = 60
	}
	errChan := make(chan error)
	evtSub := NewEventSubscriber()
	defer evtSub.Close()
	q_events := evtSub.Subscribe(sdl.K_q)
	for {
		select {
		case e := <-errChan:
			return e
		case <-q_events:
			return
			// shutdown everything
		case <-time.After(time.Second / time.Duration(fps)):
			if err = e.Render(); err != nil {
				return err
			}
		}
	}
}
