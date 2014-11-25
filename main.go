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
)

type Event struct {
	Key sdl.Keycode
	Op  func(int, int) int
}

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}

var (
	leftPaddleEvents = []Event{
		Event{sdl.K_w, sub},
		Event{sdl.K_s, add},
	}
	rightPaddleEvents = []Event{
		Event{sdl.K_UP, sub},
		Event{sdl.K_DOWN, add},
	}
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

type Paddle struct {
	PosX, PosY int
	W, H       int
	Events     []Event
}

func (p *Paddle) Update(events []sdl.Event) {
	for _, event := range events {
		switch t := event.(type) {
		case *sdl.KeyDownEvent:
			log.Print("key down event")
			for _, evt := range p.Events {
				if evt.Key == t.Keysym.Sym {
					p.PosY = evt.Op(p.PosY, 10)
				}
			}
		}
	}
}

func GetEventList() []sdl.Event {
	list := make([]sdl.Event, 10)
	for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
		list = append(list, evt)
	}
	return list
}

type Ball struct {
	W, H   float64
	DX, DY float64 // direction
	X, Y   float64
}

func detectCollision(a, b sdl.Rect) bool {
	// b.X is same height as a
	if a.X <= b.X && a.X+a.H > b.X {
		// b.Y is same height as a
		if a.Y <= b.Y && a.Y+a.W > b.Y {
			return true
		}
	}
	return false
}

func MoveBall(
	b Ball,
	leftPaddle, rightPaddle Paddle,
	arenaWidth, arenaHeight int,
) Ball {
	currentBallPos := []float64{
		float64(b.X + b.W/2),
		float64(b.Y + b.H/2),
	}
	directionVector := []float64{
		float64(b.DX),
		float64(b.DY),
	}
	// tx == currentBallPos[0] + directionVector[0] * r
	// ty == currentBallPos[1] + directionVector[1] * r

	// Let's first check for left wall
	// This makes x == 0 and y might be anything
	// 0 == currentBallPos[0] + directionVector[0] * r
	// -> -currentBallPos[0] == directionVector[0] * r
	// -> -currentBallPos[0] / directionVector[0] == r
	// if r > 1 then it's too far away and we don't care

	// We collide to the left
	if (-currentBallPos[0] / directionVector[0]) < 1 {
		// collided left wall
	}
	if (-currentBallPos[1] / directionVector[1]) < 1 {
		// collided bottom line
	}

	// Fist Fist Fist
	return Ball{
		DX:b.DX,
		DY : b.DY,
		X: b.X + b.DX,
		Y: b.Y + b.DY,
		H: b.H,
		W: b.W,
	}
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
	leftPaddle := Paddle{
		PosX:   PADDING,
		PosY:   SCREEN_HEIGHT / 2,
		Events: leftPaddleEvents,
	}
	rightPaddle := Paddle{
		PosX:   SCREEN_WIDTH - PADDING - PADDLE_WIDTH,
		PosY:   SCREEN_HEIGHT / 2,
		Events: rightPaddleEvents,
	}

	ball := Ball{
		X:  SCREEN_WIDTH / 2,
		Y:  SCREEN_HEIGHT / 2,
		DX: 0.7,
		DY: 0.3,
		W:  40,
		H:  40,
	}
	ballSurface := sdl.CreateRGBSurface(
		0,
		int32(ball.W), int32(ball.H),
		32, 0, 0, 0, 0,
	)
	ballSurface.FillRect(
		nil,
		PADDLE_COLOR,
	)
	ballTexture := renderer.CreateTextureFromSurface(ballSurface)

	// main loop
	for {
		// iterate over events

		events := GetEventList()
		leftPaddle.Update(events)
		rightPaddle.Update(events)
		ball = MoveBall(ball, leftPaddle, rightPaddle, SCREEN_WIDTH, SCREEN_HEIGHT)

		// copy whole texture to whole target
		if renderer.Copy(backgroundTexture, nil, nil) != 0 {
			return NewPongError("Copying Background Texture failed")
		}
		// TODO
		//leftPaddlePos := GetPaddlePosition(events, oldposition)
		// or something like that? and then use these information
		render := func(X, Y int) int {
			return renderer.Copy(paddleTexture, nil, &sdl.Rect{
				W: int32(PADDLE_WIDTH),
				H: int32(PADDLE_LENGTH),
				Y: int32(Y),
				X: int32(X),
			})
		}
		if render(
			leftPaddle.PosX,
			leftPaddle.PosY-PADDLE_LENGTH/2,
		) != 0 {
			return NewPongError("Copying left paddle failed")
		}
		if render(
			rightPaddle.PosX,
			rightPaddle.PosY-PADDLE_LENGTH/2,
		) != 0 {
			return NewPongError("Copying right paddle failed")
		}
		if renderer.Copy(ballTexture, nil, &sdl.Rect{
			W: int32(ball.W),
			H: int32(ball.H),
			X: int32(ball.X),
			Y: int32(ball.Y),
		}) != 0 {
			return NewPongError("Rendering ball failed")
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
