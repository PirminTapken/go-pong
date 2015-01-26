package main

import (
	//	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type SdlEventStream struct {
	closing chan chan error
	events  chan sdl.Event
}

func (ses *SdlEventStream) Close() error {
	c := make(chan error)
	ses.closing <- c
	return <-c
}

func (ses *SdlEventStream) loop() {
	var err error
	for {
		//evt := sdl.PollEvent()
		select {
		case errc := <-ses.closing:
			errc <- err
			close(ses.events)
			return
		case <-time.After(1 * time.Second):
		}
	}
}

func (ses *SdlEventStream) GetEvents() <-chan sdl.Event {
	ses.events = make(chan sdl.Event)
	ses.closing = make(chan chan error)
	go ses.loop()
	return ses.events
}

// KeyEventSubscriber sfilters key events from an event stream
type KeyEventSubscriber struct {
	eventStream <-chan sdl.Event
	closing     chan chan error
}

// NewKeyFilter returns an initialized KeyFilter
func NewKeyEventSubscriber(eventStream <-chan sdl.Event) *KeyEventSubscriber {
	kes := &KeyEventSubscriber{eventStream, make(chan chan error)}
	go kes.loop()
	return kes
}

func (kes *KeyEventSubscriber) loop() {
	var err error
	for {
		select {
		case errc := <-kes.closing:
			errc <- err
			return
		case <-kes.eventStream:
			// distribute
		}
	}
}

func (kes *KeyEventSubscriber) KeyEvents(sdl.Keycode) <-chan sdl.Event {
	// FIXME
	return kes.eventStream
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
