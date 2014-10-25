package main

import (
	"testing"
)

func TestPlainSDLError(t *testing.T) {
	err := NewSdlError("wurst")
	if err.Error() != "wurst" {
		t.Errorf("%v != %v", err.Error(), "wurst")
	}
}

func TestRealSDLError(t *testing.T) {
	// TODO:
	// We need to create a real sdl error
	// then create a custom sdl error and see
	// what happens
	t.SkipNow()
}
