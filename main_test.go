package main

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
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

// We want to know if the background is really filled with the background color
func TestBackgroundFill(t *testing.T) {
	size := struct {
		W int
		H int
	}{
		10,
		10,
	}
	window, renderer := sdl.CreateWindowAndRenderer(size.W, size.H, 0)
	err := fillBackground(window, renderer)
	if err != nil {
		t.Error("fill background error", err)
	}

	surface := window.GetSurface()
	pixels := surface.Pixels()
	firstPixel := uint32(pixels[0])<<8 + uint32(pixels[1])<<8 + uint32(pixels[2])<<8 + uint32(pixels[3])

	if firstPixel != uint32(BACKGROUND) {
		t.Error("FirstPixel !=", BACKGROUND)
	}

	lastPixel := uint32(pixels[surface.PixelNum()-3])<<8 + uint32(pixels[surface.PixelNum()-2])<<8 + uint32(pixels[surface.PixelNum()-1])<<8 + uint32(pixels[surface.PixelNum()])

	if lastPixel != uint32(BACKGROUND) {
		t.Error("LastPixel !=", BACKGROUND)
	}
}

func TestNewPaddel(t *testing.T) {
	paddel, err := NewPaddel(0, 0)
	if err != nil {
		t.Error("error creating paddel", err)
	}
	if paddel == nil {
		t.Error("paddel is nil for some reason")
	}
}

func TestPaddleGetCoord(t *testing.T) {
	paddel, err := NewPaddel(0, 0)
	if err != nil {
		t.Error("error creating paddel", err)
	}
	if paddel == nil {
		t.Error("paddel is nil for some reason")
	}
	coord, err := paddel.GetCoord()
	if err != nil {
		t.Error("Paddel.GetCoord returned error", err)
	}
	if coord.X != 0 || coord.Y != 0 {
		t.Error("Coordinates wrong!", coord)
	}
}
