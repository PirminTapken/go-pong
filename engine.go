package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/mobile/sprite"
	"sync"
)

const (
	BACKGROUND_COLOR = uint32(0x000000)
	PADDLE_LENGTH    = 0.2
	PADDLE_WIDTH     = 0.05
	PADDLE_COLOR     = uint32(0xffffff)
)

// A Scene Graph, also partially stolen from golang.org/x/mobile/sprite
// but safe for concurrency
type Node struct {
	Arranger                                                sprite.Arranger
	SubTex                                                  sprite.SubTex
	LastChild, FirstChild, NextSibling, PrevSibling, Parent *Node
	M                                                       sync.RWMutex
}

// AppendChild appends a node to the graph
func (n *Node) AppendChild(c *Node) {
	n.M.Lock()
	defer n.M.Unlock()
	c.M.Lock()
	defer c.M.Unlock()
	if c.Parent != nil || c.PrevSibling != nil || c.NextSibling != nil {
		panic("Node: Child already somewhere in tree")
	}

	// taken from golang.org/x/mobile/sprite
	last := n.LastChild
	if last != nil {
		last.M.Lock()
		last.NextSibling = c
		last.M.Unlock()
	} else {
		n.FirstChild = c
	}
	c.PrevSibling = last
	c.Parent = n
	n.LastChild = c
}

// RemoveChild removes child from tree
func (n *Node) RemoveChild(c *Node) {
	n.M.Lock()
	defer n.M.Unlock()
	c.M.Lock()
	defer c.M.Unlock()

	prevSibling := c.PrevSibling
	nextSibling := c.NextSibling
	if prevSibling != nil {
		prevSibling.M.Lock()
		prevSibling.NextSibling = nextSibling
		prevSibling.M.Unlock()
	}
	if nextSibling != nil {
		nextSibling.M.Lock()
		nextSibling.PrevSibling = prevSibling
		nextSibling.M.Unlock()
	}

	if n.FirstChild == c {
		n.FirstChild = c.NextSibling
	}
	if n.LastChild == c {
		n.LastChild = c.PrevSibling
	}

	c.Parent = nil
	c.NextSibling = nil
	c.PrevSibling = nil
}

// SDLEngine is an 2d engine using SDL
type SDLEngine struct {
	cleanupFns []func()
	window     *sdl.Window
	renderer   *sdl.Renderer
	// Thread is public so it can be used by other
	// goroutines that need to do stuff in the sdl
	// thread
	Thread *Thread
}

// NewSDLEngine creates the engine.
// This basically creates the background texture and stores it away
func NewSDLEngine(windowName string, X, Y, W, H int) (e *SDLEngine, err error) {
	e = &SDLEngine{Thread: NewThread(),
		cleanupFns: make([]func(), 0)}
	err = e.sdlInit()
	if err != nil {
		return e, err
	}
	err = e.CreateWindowAndRenderer(W, H, 0)
	e.SetTitle(windowName)
	return e, err
}

// Close closes the engine
// error is always nil and just there to match
// io.Closer
func (e *SDLEngine) Close() error {
	e.cleanup()
	err := e.Thread.Close()
	return err
}

func (e *SDLEngine) Title() string {
	return e.Thread.Exec(func() interface{} {
		return e.window.GetTitle()
	}).(string)
}

func (e *SDLEngine) SetTitle(s string) {
	_ = e.Thread.Exec(func() interface{} {
		e.window.SetTitle(s)
		return nil
	})
}

// init calls sdl init in sdl thread
func (e *SDLEngine) sdlInit() error {
	r := e.Thread.Exec(func() interface{} {
		// This is neccessary otherwise nil error can't be converted to
		// interface and back somehow...
		e := struct{ err error }{}
		if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
			e.err = sdl.GetError()
		}
		return e
	}).(struct{ err error })
	e.cleanupFns = append(e.cleanupFns, sdl.Quit)
	return r.err
}

// cleanup cleans everything up
func (e *SDLEngine) cleanup() {
	for i := len(e.cleanupFns); i > 0; i-- {
		e.Thread.Exec(func() interface{} {
			e.cleanupFns[i-1]()
			return nil
		})
	}
}

func (e *SDLEngine) CreateWindowAndRenderer(w, h int, flags uint32) error {
	type resp struct {
		w *sdl.Window
		r *sdl.Renderer
		e error
	}
	r := e.Thread.Exec(func() interface{} {
		w, r, e := sdl.CreateWindowAndRenderer(w, h, flags)
		return resp{w: w, r: r, e: e}
	}).(resp)
	if r.e != nil {
		return r.e
	}
	e.window = r.w
	e.renderer = r.r
	e.cleanupFns = append(e.cleanupFns, e.window.Destroy)
	e.cleanupFns = append(e.cleanupFns, e.renderer.Destroy)
	return nil
}
