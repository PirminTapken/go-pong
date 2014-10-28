package main

import (
	"log"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDEL_LENGTH    = 100
	PADDEL_WIDTH     = 20
	PADDEL_COLOR     = uint32(0xffffff)
	SCREEN_WIDTH     = 640
	SCREEN_HEIGHT    = 480
	NAME             = "Pong"
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

// run on the renderer
func Run(
	screen struct {
		W int
		H int
	},
	renderer *sdl.Renderer,
) error {
	background := sdl.CreateRGBSurface(
		0,
		int32(screen.W),
		int32(screen.H),
		32,
		0, 0, 0, 0,
	)
	if background == nil {
		return NewPongError("background creation failed")
	}
	if background.FillRect(
		&sdl.Rect{
			X: 0,
			Y: 0,
			W: int32(screen.W),
			H: int32(screen.H),
		},
		BACKGROUND_COLOR,
	) != 0 {
		return NewPongError("Filling background failed")
	}
	backgroundTexture := renderer.CreateTextureFromSurface(background)
	if backgroundTexture == nil {
		return NewPongError("Creating Background Texture failed")
	}
	// copy whole texture to whole target
	if renderer.Copy(backgroundTexture, nil, nil) != 0 {
		return NewPongError("Copying Background Texture failed")
	}
	return nil
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
	err := Run(struct {
		W int
		H int
	}{SCREEN_WIDTH, SCREEN_HEIGHT}, renderer)
	if err != nil {
		log.Fatal(err)
	}
}
