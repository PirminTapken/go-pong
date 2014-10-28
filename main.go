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
	if background.FillRect(nil, BACKGROUND_COLOR) != 0 {
		return NewPongError("Filling background failed")
	}
	backgroundTexture := renderer.CreateTextureFromSurface(background)
	if backgroundTexture == nil {
		return NewPongError("Creating Background Texture failed")
	}
	paddleSurface := sdl.CreateRGBSurface(0, PADDLE_WIDTH, PADDLE_LENGTH, 32, 0, 0, 0, 0)
	if paddleSurface == nil {
		return NewPongError("Creating paddle surface failed")
	}
	if paddleSurface.FillRect(nil, PADDLE_COLOR) != 0 {
		return NewPongError("Filling Paddle Surface failed")
	}
	paddleTexture := renderer.CreateTextureFromSurface(paddleSurface)
	if paddleTexture == nil {
		return NewPongError("Creating paddle texture failed")
	}
	leftPaddlePos := SCREEN_HEIGHT / 2
	rightPaddlePos := SCREEN_HEIGHT / 2
	// main loop
	var evtQueue []string
	for {
		evtQueue = make([]string, 5)
		// iterate over events
		for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			log.Print("polling events")
			switch t := evt.(type) {
			case *sdl.QuitEvent:
				return NewPongError("Quitting the hard way ;)")
			case *sdl.KeyDownEvent:
				if t.Keysym.Sym == sdl.K_UP {
					evtQueue = append(evtQueue, "upkey")
				}
				if t.Keysym.Sym == sdl.K_DOWN {
					evtQueue = append(evtQueue, "downkey")
				}
			}
		}

		log.Print("processing events")
		for _, evt := range evtQueue {
			if evt == "upkey" {
				rightPaddlePos -= 10
			}
			if evt == "downkey" {
				rightPaddlePos += 10
			}
		}
		// copy whole texture to whole target
		if renderer.Copy(backgroundTexture, nil, nil) != 0 {
			return NewPongError("Copying Background Texture failed")
		}
		if renderer.Copy(paddleTexture, nil, &sdl.Rect{
			W: PADDLE_WIDTH,
			H: PADDLE_LENGTH,
			Y: int32(leftPaddlePos - PADDLE_LENGTH/2),
			X: 10,
		}) != 0 {
			return NewPongError("Copying left paddle failed")
		}
		if renderer.Copy(paddleTexture, nil, &sdl.Rect{
			W: PADDLE_WIDTH,
			H: PADDLE_LENGTH,
			Y: int32(rightPaddlePos - PADDLE_LENGTH/2),
			X: SCREEN_WIDTH - 10 - PADDLE_WIDTH,
		}) != 0 {
			return NewPongError("Copying left paddle failed")
		}
		renderer.Present()
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
