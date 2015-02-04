package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
)

// TestCloseEventStream tests if closing our event stream works
func TestCloseEventStream(t *testing.T) {
	es := &SdlEventStream{}
	events := es.Receive()
	t.Log("Opened channel, calling close in background")
	err := es.Close()
	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
	if evt := <-events; evt != nil {
		t.Errorf("events returned something other than it's zero value nil: %v", evt)
	}
}

func TestCloseEventSubscriber(t *testing.T) {
	es := NewEventSubscriber()
	if err := es.Close(); err != nil {
		t.Errorf("closing returned error: %v", err)
	}

}

func TestReceiveSpecificKeyEvent(t *testing.T) {
	es := NewEventSubscriber()
	evtChan := es.Subscribe(sdl.K_UP)
	fakeEvent := &sdl.KeyDownEvent{
		Keysym: sdl.Keysym{
			Sym:      sdl.K_UP,
			Scancode: sdl.SCANCODE_UP,
			Mod:      sdl.KMOD_NONE,
		},
		Type:  sdl.KEYDOWN,
		State: sdl.PRESSED,
	}
	sdl.PushEvent(fakeEvent)
	evt := <-evtChan
	switch e := evt.(type) {
	case *sdl.KeyDownEvent:
		if e.Keysym.Sym != sdl.K_UP {
			t.Errorf("Key is %v instead of %v", e.Keysym.Sym, sdl.K_UP)
		}
	default:
		t.Errorf("Event was %v instead of KeyDownEvent", evt)
	}
}
