package main

import (
	"log"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	BACKGROUND    = uint32(0x000000)
	PADDEL_LENGTH = 100
	PADDEL_WIDTH  = 20
	PADDEL_COLOR  = uint32(0xffffff)
)

type SdlError struct {
	msg      string
	sdlError error
}

func NewSdlError(msg string) *SdlError {
	return &SdlError{
		msg:      msg,
		sdlError: sdl.GetError(),
	}
}

func (e *SdlError) Error() string {
	return strings.Join(
		[]string{
			e.msg,
			e.sdlError.Error(),
		},
		": ",
	)
}

func fillBackground(window *sdl.Window, renderer *sdl.Renderer) error {
	w, h := window.GetSize()
	background := sdl.CreateRGBSurface(
		0,
		int32(w),
		int32(h),
		32,
		0, 0, 0, 0,
	)
	if background == nil {
		return NewSdlError("background nil")
	}
	allRect := &sdl.Rect{
		W: int32(w),
		H: int32(h),
		X: 0,
		Y: 0,
	}
	if background.FillRect(
		allRect,
		BACKGROUND,
	) != 0 {
		return NewSdlError("Background fill failed")
	}
	tex := renderer.CreateTextureFromSurface(background)
	if tex == nil {
		return NewSdlError("CreateTextureFromSurfaceFailed")
	}
	background.Free()
	if renderer.Copy(tex, nil, allRect) != 0 {
		return NewSdlError("RenderCopy Failed")
	}
	tex.Destroy()
	return nil
}

type EventListener interface {
	PutEvent(event sdl.Event)
}

// Right now this type needs to be created manually
// I am too lazy now to write a registrar for new
// listener
type EventDistributor struct {
	Listener []EventListener
}

// let this event distributor run in own goroutine
func (ed *EventDistributor) Run() {
	var event sdl.Event
	for {
		event = sdl.PollEvent()
		if event != nil {
			for _, listener := range ed.Listener {
				go listener.PutEvent(event)
			}
		}
	}
}

type Coord struct {
	X int
	Y int
}

type Paddel struct {
	coord        Coord
	surface      *sdl.Surface
	errChan      chan error
	coordChan    chan Coord
	newCoordChan chan Coord
}

func NewPaddel(initialX, initialY int) (*Paddel, error) {
	p := Paddel{
		coord:        Coord{X: initialX, Y: initialY},
		surface:      sdl.CreateRGBSurface(0, int32(PADDEL_WIDTH), int32(PADDEL_LENGTH), 32, 0, 0, 0, 0),
		errChan:      make(chan error),
		coordChan:    make(chan Coord),
		newCoordChan: make(chan Coord),
	}
	if p.surface == nil {
		return nil, NewSdlError("Creating Surface for new Paddel failed")
	}
	if p.surface.FillRect(
		&sdl.Rect{
			X: 0, Y: 0,
			W: PADDEL_WIDTH,
			H: PADDEL_LENGTH,
		},
		PADDEL_COLOR,
	) != 0 {
		return nil, NewSdlError("NewPaddel Surface filling failed")
	}
	go p.run()
	return &p, nil
}

func (p *Paddel) PutEvent(evt sdl.Event) {
	log.Print("received event", evt)
	switch t := evt.(type) {
	case sdl.CommonEvent:
		log.Print("key down event")
		switch t.Keysym.Sym {
		case sdl.K_UP:
			log.Print("up event")
			coord, _ := p.GetCoord()
			coord.Y--
			p.newCoordChan <- coord
		case sdl.K_DOWN:
			log.Print("down event")
			coord, _ := p.GetCoord()
			coord.Y++
			p.newCoordChan <- coord
		}
	}
}

func (p *Paddel) GetCoord() (Coord, error) {
	return <-p.coordChan, nil
}

func (p *Paddel) Render(renderer *sdl.Renderer) error {
	tex := renderer.CreateTextureFromSurface(p.surface)
	if tex == nil {
		return NewSdlError("Creating Texture for Paddel render failed")
	}
	coord, err := p.GetCoord()
	if err != nil {
		return err
	}
	if renderer.Copy(tex, nil, &sdl.Rect{
		X: int32(coord.X),
		Y: int32(coord.Y),
		W: PADDEL_WIDTH,
		H: PADDEL_LENGTH,
	}) != 0 {
		return NewSdlError("Render Copy for Paddel failed")
	}
	return nil
}

// never call this!
func (p *Paddel) setCoords(newCoords Coord) {
	p.coord = newCoords
}

// internal runner
func (p *Paddel) run() {
	for {
		select {
		case newCoords := <-p.newCoordChan:
			p.setCoords(newCoords)
		case p.coordChan <- p.coord:
			// We just wanted to answer a request
		default:
			continue
		}
	}
}

// Run the program after all init stuff
func Run(window *sdl.Window, renderer *sdl.Renderer) error {
	log.Print("starting run")
	var err error
	w, h := window.GetSize()
	leftPaddel, err := NewPaddel(
		10,
		(h-PADDEL_LENGTH)/2,
	)
	if err != nil {
		return err
	}
	rightPaddel, err := NewPaddel(
		w-10-PADDEL_WIDTH,
		(h-PADDEL_LENGTH)/2,
	)
	if err != nil {
		return err
	}

	eventDistributor := &EventDistributor{
		Listener: []EventListener{leftPaddel, rightPaddel},
	}

	go eventDistributor.Run()

	// main loop
	for {
		err = fillBackground(window, renderer)
		if err != nil {
			return err
		}
		for _, p := range []*Paddel{leftPaddel, rightPaddel} {
			err = p.Render(renderer)
			if err != nil {
				return err
			}
		}
		renderer.Present()
	}
	return nil
}

func main() {
	log.Print("Hello World!")
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()
	window := sdl.CreateWindow(
		"Pong",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		640,
		480,
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
	err := Run(window, renderer)
	if err != nil {
		log.Fatal(err)
	}
}
