package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"strconv"
	"time"
)

type registerMessage struct {
	Response    chan int
	InitialDest *sdl.Rect
}

type updateMessage struct {
	Response chan bool
	SpriteId int
	Rect     *sdl.Rect
}

// Engine is our little nice graphics engine
type Engine struct {
	Renderer       *sdl.Renderer
	Background     *sdl.Texture
	ObjectTexture  *sdl.Texture
	RegisterSprite chan registerMessage
	UpdateSprite   chan updateMessage
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

type Sprite struct {
	Dest *sdl.Rect
}

// Render the world
func (e *Engine) Render(world []Sprite) (err error) {
	err = e.Renderer.Copy(e.Background, nil, nil)
	if err != nil {
		return err
	}
	for _, item := range world {
		e.Renderer.Copy(item.Tex, nil, item.Dest)
	}
	e.Renderer.Present()
	return nil
}

// Run is the loop of the engine
// it calls the render func
func (e *Engine) Run() (err error) {
	fps, err := strconv.Atoi(os.Getenv("FPS"))
	if err != nil {
		fps = 60
	}
	errChan := make(chan error)
	evtSub := NewEventSubscriber()
	defer evtSub.Close()
	q_events := evtSub.Subscribe(sdl.K_q)

	var world []Sprite

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
