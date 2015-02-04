package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"testing"
	"time"
)

// TestCloseEventStream tests if closing our event stream works
func TestCloseEventStream(t *testing.T) {
	es := NewSdlEventStream(NewThread())
	if es == nil {
		t.Error("SdlEventStream is nil")
	}
	events := es.Receive()
	if events == nil {
		t.Error("events channel is nil")
	}
	t.Log("Opened channel, calling close in background")
	errc := make(chan error)
	// do it in background to give it time
	go func() {
		errc <- es.Close()
	}()
	select {
	case err := <-errc:
		if err != nil {
			t.Errorf("An error occured: %v", err)
		}
		// we have finished in time
	case <-time.After(3 * time.Second):
		t.Fatal("timeout!")
	}
	if evt := <-events; evt != nil {
		t.Errorf("events returned something other than it's zero value nil: %v", evt)
	}
}

func TestCloseEventSubscriber(t *testing.T) {
	es := NewEventSubscriber()
	errc := make(chan error)
	go func() {
		errc <- es.Close()
	}()
	select {
	case err := <-errc:
		if err != nil {
			t.Errorf("closing returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timeout!")
	}
}

func TestReceiveSpecificKeyEvent(t *testing.T) {
	done := make(chan bool)
	go func() {
		es := NewEventSubscriber()
		t.Log("Subscribe")
		evtChan := es.Subscribe(sdl.K_UP)
		t.Log("Create Fake Event")
		fakeEvent := &sdl.KeyDownEvent{
			Keysym: sdl.Keysym{
				Sym:      sdl.K_UP,
				Scancode: sdl.SCANCODE_UP,
				Mod:      sdl.KMOD_NONE,
				Unicode:  0,
			},
			Type:      sdl.KEYDOWN,
			Timestamp: 0,
			WindowID:  0,
			State:     sdl.PRESSED,
			Repeat:    0,
		}
		t.Log("Push fake Event")
		sdl.PushEvent(fakeEvent)
		select {
		case evt := <-evtChan:
			switch e := evt.(type) {
			case *sdl.KeyDownEvent:
				if e.Keysym.Sym != sdl.K_UP {
					t.Errorf("Key is %v instead of %v", e.Keysym.Sym, sdl.K_UP)
				}
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout after 5 seconds")
		}
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Error("Timeout Nothing went well")
	}
}
