package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"strconv"
	"time"
)

// Engine is our little nice graphics engine
type Engine struct {
	Renderer      *sdl.Renderer
	Background    *sdl.Texture
	ObjectTexture *sdl.Texture
}

// CreateEngine creates the engine.
// This basically creates the background texture and stores it away
func CreateEngine(renderer *sdl.Renderer) (*Engine, error) {
	bgSurface := sdl.CreateRGBSurface(
		0, SCREEN_WIDTH, SCREEN_HEIGHT, 32,
		0, 0, 0, 0,
	)
	if bgSurface == nil {
		return nil, fmt.Errorf("surface creation failed: %v", sdl.GetError())
	}
	if bgSurface.FillRect(nil, BACKGROUND_COLOR) != 0 {
		return nil, fmt.Errorf("filling background failed: %v", sdl.GetError())
	}
	backgroundTexture, err := renderer.CreateTextureFromSurface(bgSurface)
	if err != nil {
		return nil, fmt.Errorf("Creating BackgroundTexture failed:", err)
	}
	objectSurface := sdl.CreateRGBSurface(0, 1, 1, 32, 0, 0, 0, 0) // just a single color
	objectSurface.FillRect(nil, PADDLE_COLOR)
	objectTexture, err := renderer.CreateTextureFromSurface(objectSurface)
	if err != nil {
		return nil, fmt.Errorf("Creating ObjectTexture failed:", err)
	}
	engine := &Engine{
		Renderer:      renderer,
		Background:    backgroundTexture,
		ObjectTexture: objectTexture,
	}
	return engine, nil
}

// Render the world
func (e *Engine) Render() (err error) {
	e.Renderer.Copy(e.Background, nil, nil)
	if err != nil {
		return err
	}
	// TODO what's on with the rest?
	e.Renderer.Present()
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
	for {
		select {
		case e := <-errChan:
			return e
		case <-time.After(time.Second / time.Duration(fps)):
			if err = e.Render(); err != nil {
				return err
			}
		}
	}
}
