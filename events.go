package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

// EventSpreader spreads sdl.Events from a channel to several listening channels
// It synchronizes so each event needs to be recieved by each channel until it continues
func EventSpreader(inChan <-chan sdl.Event, outChan ...*EventHandler) {
	for {
		wg := new(sync.WaitGroup)
		wg.Add(len(outChan))
		e := <-inChan
		for _, oc := range outChan {
			go func() {
				oc.EventChan <- e
				wg.Done()
			}()
		}
	}
}

// PollSdlEvents recieves sdl events and puts them on a channel
// which is returned upon calling this function
// also it returns a quit channel which can be used to quit GetSdlEvents.
func PollSdlEvents(quit chan bool) (sdlEvent <-chan sdl.Event) {
	c := make(chan sdl.Event)
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				c <- sdl.PollEvent()
			}
		}
	}()
	return c
}

// EventHandler should not be created directly but using it's functions
// As it contains only channels it's safe though
type EventHandler struct {
	Quit      chan bool
	EventChan chan sdl.Event
}

// Returns an initialized event handler
func NewEventHandler() *EventHandler {
	return &EventHandler{
		EventChan: make(chan sdl.Event),
	}
}

// KeyDownEventHandler tackles KeyDown Events and executes action on correct key
func KeyDownEventHandler(key sdl.Keycode, action func(), quit chan bool) chan sdl.Event {
	eventChan := make(chan sdl.Event)
	go func() {
		for {
			select {
			case e := <-eventChan:
				switch et := e.(type) {
				case *sdl.KeyDownEvent:
					switch et.Keysym.Sym {
					case key:
						action()
					}
				}
			case <-quit:
				quit <- true
				return
			}
		}
	}()
	return eventChan
}

// LoopEvents Takes care of the event handling
// returns a channel that emits a signal when a quit
// event is received
func LoopEvents(errChan chan error, quit chan bool) chan bool {

	quitPollingEvents := make(chan bool)
	sdlEvents := PollSdlEvents(quitPollingEvents)

	quitEvent := make(chan bool)
	keyQChan := KeyDownEventHandler(sdl.K_q,
		func() {
			quitEvent <- true
		},
		quit)
	return quitEvent
}
