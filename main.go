package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDLE_LENGTH    = 0.2
	PADDLE_WIDTH     = 0.05
	PADDLE_COLOR     = uint32(0xffffff)
	SCREEN_WIDTH     = 640
	SCREEN_HEIGHT    = 480
	NAME             = "Pong"
	PADDING          = 10
	VERSION          = "0.1"
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

type Object struct {
	W, H, X, Y float64
	DX, DY     float64
}

func GetEventList() []sdl.Event {
	list := make([]sdl.Event, 10)
	for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
		list = append(list, evt)
	}
	return list
}

type Engine struct {
	Renderer   *sdl.Renderer
	Background *sdl.Texture
}

// Render the world
func (e *Engine) Render(universeBus chan map[string]Object) (err error) {
	log.Print("Get unierse to render")
	universe := <-universeBus
	defer func() {
		log.Print("Free universe after render")
		universeBus <- universe
	}()
	e.Renderer.Copy(e.Background, nil, nil)
	s := sdl.CreateRGBSurface(0, 1, 1, 32, 0, 0, 0, 0) // just a single color
	s.FillRect(nil, PADDLE_COLOR)
	t, err := e.Renderer.CreateTextureFromSurface(s)
	if err != nil {
		return err
	}
	for _, obj := range universe {
		err = e.Renderer.Copy(
			t, nil,
			&sdl.Rect{
				X: int32((obj.X - obj.W/2) * SCREEN_WIDTH),
				Y: int32((obj.Y - obj.H/2) * SCREEN_HEIGHT),
				H: int32(obj.H * SCREEN_HEIGHT),
				W: int32(obj.W * SCREEN_WIDTH),
			},
		)
		if err != nil {
			return err
		}
	}
	e.Renderer.Present()
	return nil
}

type Direction bool

const (
	UP   Direction = true
	DOWN Direction = false
)

// UpdatePaddle updates a paddle
func UpdatePaddle(universeBus chan map[string]Object, errChan chan error, paddle string, d Direction) {
	v := 0.0
	switch d {
	case UP:
		v = v - 0.01
	case DOWN:
		v = v + 0.01
	}
	u := <-universeBus
	go func() {
		universeBus <- u
	}()
	tmp, ok := u[paddle]
	if !ok {
		errChan <- fmt.Errorf(`Key "%s" does not exist in our universe!`, paddle)
		return
	}
	if 0 < tmp.Y+v && tmp.Y+v < 1 {
		tmp.Y = tmp.Y + v
	}
	// assign updated paddle back as we don't have a pointer (yet)
	u[paddle] = tmp
}

type WallIntersection struct {
	IntersectAt float64
	Wall        *Line
}

// WallIntersections returns wall intersections
func WallIntersections(walls []*Line, line *Line) []WallIntersection {
	intersections := make([]WallIntersection, 4)
	for _, wall := range walls {
		log.Print("wall: ", wall)
		var intersection WallIntersection
		intersection.IntersectAt = line.Intersect(wall)
		intersection.Wall = wall
	}
	return intersections
}

func UpdateBall(universeBus chan map[string]Object, errChan chan error, d time.Duration) {
	walls := []*Line{
		&Line{&Vector2d{0, 0}, &Vector2d{0, 1}},
		&Line{&Vector2d{0, 1}, &Vector2d{1, 1}},
		&Line{&Vector2d{1, 1}, &Vector2d{1, 0}},
		&Line{&Vector2d{1, 0}, &Vector2d{0, 0}},
	}

	u := <-universeBus
	defer func() {
		universeBus <- u
	}()
	ball := u["Ball"]
	defer func() {
		u["Ball"] = ball
	}()

	initialDir := &Vector2d{ball.DX * d.Seconds(), ball.DY * d.Seconds()}
	initialPos := &Vector2d{ball.X, ball.Y}

	pos := initialPos.Copy()
	dir := initialDir.Copy()
	newPos := pos.Add(dir)
	line := &Line{pos, newPos}
	for line.Vector2d().Len() > 0 {

		wallmap := make(map[float64]*Line, 4)

		h := make([]float64, 4)

		for i, wall := range walls {
			h[i] = line.Intersect(wall)
			wallmap[h[i]] = wall
		}

		sort.Float64s(h)

		c := 0
		for _, f := range h {
			if 0 < f && f < 1 {
				// we do intersect
				hitPos := pos.Add(dir.Scale(f))
				remainder := line.Vector2d().Len() - f
				dir = wallmap[f].Vector2d().Reflect(dir.Scale(remainder))
				newPos = hitPos.Add(dir)
				pos = hitPos
				line = &Line{pos, newPos}
				c = c + 1
			}
		}
		if c == 0 {
			newPos = pos.Add(dir)
			line = &Line{pos, newPos}
			break
		}
	}
	ball.X = newPos[0]
	ball.Y = newPos[1]
}

func LoopEvents(universeBus chan map[string]Object, errChan chan error, quit chan bool) {
	eventList := GetEventList()
	for _, event := range eventList {
		switch e := event.(type) {
		case *sdl.KeyDownEvent:
			switch e.Keysym.Sym {
			case sdl.K_q:
				quit <- true
			case sdl.K_DOWN:
				go UpdatePaddle(universeBus, errChan, "Right Paddle", DOWN)
			case sdl.K_UP:
				go UpdatePaddle(universeBus, errChan, "Right Paddle", UP)
			case sdl.K_w:
				go UpdatePaddle(universeBus, errChan, "Left Paddle", UP)
			case sdl.K_s:
				go UpdatePaddle(universeBus, errChan, "Left Paddle", DOWN)
			}

		}
	}
}

// Run the game
func Run(e *Engine, universeBus chan map[string]Object) (err error) {
	fps, err := strconv.Atoi(os.Getenv("FPS"))
	if err != nil {
		fps = 60
	}
	clockChan := time.Tick(time.Second / time.Duration(fps*int(time.Second)))
	quit := make(chan bool, 1)
	errChan := make(chan error)
	now := time.Now()
	last := now
	for {
		go UpdateBall(universeBus, errChan, now.Sub(last))
		go LoopEvents(universeBus, errChan, quit)
		select {
		case <-quit:
			return nil
		case e := <-errChan:
			return e
		default:
			// continue
		}
		if err = e.Render(universeBus); err != nil {
			return err
		}
		// wait for tick
		last = now
		now = <-clockChan
	}
}

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow(
		NAME,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
		0,
	)
	if err != nil {
		log.Fatal("creating window failed", err)
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(
		window,
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		log.Fatal("creating renderer failed", err)
	}
	defer renderer.Destroy()
	bgSurface := sdl.CreateRGBSurface(
		0, SCREEN_WIDTH, SCREEN_HEIGHT, 32,
		0, 0, 0, 0,
	)
	if bgSurface == nil {
		log.Fatal("surface creation failed", sdl.GetError())
	}
	if bgSurface.FillRect(nil, BACKGROUND_COLOR) != 0 {
		log.Fatal("filling background failed", sdl.GetError())
	}
	backgroundTexture, err := renderer.CreateTextureFromSurface(bgSurface)
	if err != nil {
		log.Fatal("Creating BackgroundTexture failed:", err)
	}
	engine := &Engine{
		Renderer:   renderer,
		Background: backgroundTexture,
	}
	universe := map[string]Object{
		"Left Paddle": Object{
			W: PADDLE_WIDTH,
			H: PADDLE_LENGTH,
			X: 0.1,
			Y: 0.5,
		},
		"Right Paddle": Object{
			W: PADDLE_WIDTH,
			H: PADDLE_LENGTH,
			X: 0.9,
			Y: 0.5,
		},
		"Ball": Object{
			W:  0.1,
			H:  0.1,
			X:  0.5,
			Y:  0.5,
			DX: 0.2,
			DY: 0.1,
		},
	}
	universeBus := make(chan map[string]Object, 1)
	universeBus <- universe
	err = Run(engine, universeBus)
	if err != nil {
		log.Fatal(err)
	}
}
