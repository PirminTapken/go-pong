package main

import (
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
		evt := sdl.PollEvent()
		select {
		case errc := <-ses.closing:
			errc <- err
			close(ses.events)
			return
		case <-time.After(1 * time.Second):
		default:
			if evt != nil {
				ses.events <- evt
			}
		}
	}
}

// Receive returns a read only channel where sdl Events will
// be submitted
func (ses *SdlEventStream) Receive() <-chan sdl.Event {
	// This makes it possible to call this more than once
	// but don't have this starting two loops or more
	if ses.events == nil {
		ses.events = make(chan sdl.Event)
		ses.closing = make(chan chan error)
		go ses.loop()
	}
	return ses.events
}

type subscribeMessage struct {
	Key  sdl.Keycode
	Chan chan sdl.Event
}

// KeyEventSubscriber sfilters key events from an event stream
type EventSubscriber struct {
	eventStream   *SdlEventStream
	eventChan     <-chan sdl.Event
	closing       chan chan error
	recepients    map[sdl.Keycode][]chan sdl.Event
	newSubscriber chan subscribeMessage
}

// NewKeyFilter returns an initialized KeyFilter
func NewEventSubscriber() *EventSubscriber {
	es := &EventSubscriber{}
	es.eventStream = &SdlEventStream{}
	es.eventChan = es.eventStream.Receive()
	es.recepients = make(map[sdl.Keycode][]chan sdl.Event)
	es.closing = make(chan chan error)
	es.newSubscriber = make(chan subscribeMessage)
	go es.loop()
	return es
}

func (es *EventSubscriber) loop() {
	var err error
	for {
		select {
		case errc := <-es.closing:
			errc <- err
			return
		case evt := <-es.eventChan:
			switch et := evt.(type) {
			case *sdl.KeyDownEvent:
				for _, c := range es.recepients[et.Keysym.Sym] {
					c <- evt
				}
			}
		case sub := <-es.newSubscriber:
			if es.recepients[sub.Key] == nil {
				es.recepients[sub.Key] = []chan sdl.Event{sub.Chan}
			} else {
				es.recepients[sub.Key] = append(es.recepients[sub.Key], sub.Chan)
			}
		}
	}
}

func (es *EventSubscriber) Subscribe(key sdl.Keycode) <-chan sdl.Event {
	c := make(chan sdl.Event)
	es.newSubscriber <- subscribeMessage{key, c}
	return c
}

func (es *EventSubscriber) Close() error {
	errc := make(chan error)
	es.closing <- errc
	return <-errc
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
