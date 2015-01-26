package main

import (
	//	"github.com/veandco/go-sdl2/sdl"
	"testing"
	"time"
)

// TestCloseEventStream tests if closing our event stream works
func TestCloseEventStream(t *testing.T) {
	es := &SdlEventStream{}
	events := es.GetEvents()
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

func TestReceiveSpecificKeyEvent(t *testing.T) {
	t.Error("not implemented")
}
