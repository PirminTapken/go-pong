package main

import (
	"log"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDLE_LENGTH    = 100
	PADDLE_WIDTH     = 20
	PADDLE_COLOR     = uint32(0xffffff)
	SCREEN_WIDTH     = 640
	SCREEN_HEIGHT    = 480
	NAME             = "Pong"
	PADDING          = 10
	VERSION = "0.1"
)

type PongError struct {
	SDLError error
	Msg      string
}

func NewPongError(msg string) *PongError {
	e := &PongError{
		SDLError: sdl.GetError(),
		Msg:      msg,
	}
	return e
}

func (e *PongError) Error() string {
	// we don't need do much if we don't have sdl errors present
	if e.SDLError == nil {
		return e.Msg
	}
	errMsg := strings.Join(
		[]string{
			e.Msg,
			e.SDLError.Error(),
		},
		": ",
	)
	return errMsg
}


func GetEventList() []sdl.Event {
	list := make([]sdl.Event, 10)
	for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
		list = append(list, evt)
	}
	return list
}

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()
	window := sdl.CreateWindow(
		NAME,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
		0,
	)
	if window == nil {
		log.Fatal(sdl.GetError(), "creating window failed")
	}
	defer window.Destroy()
	renderer := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if renderer == nil {
		log.Fatal(sdl.GetError(), "creating renderer failed")
	}
	defer renderer.Destroy()
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}
