package main

import (
	"testing"
	"time"

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

func getPaddelForTest(t *testing.T) *Paddel {
	paddel, err := NewPaddel(0, 0)
	if err != nil {
		t.Error("error creating paddel", err)
	}
	if paddel == nil {
		t.Error("paddel is nil for some reason")
	}
	return paddel
}

func TestNewPaddel(t *testing.T) {
	getPaddelForTest(t)
}

func TestPaddleGetCoord(t *testing.T) {
	paddel := getPaddelForTest(t)
	coord, err := paddel.GetCoord()
	if err != nil {
		t.Error("Paddel.GetCoord returned error", err)
	}
	if coord.X != 0 || coord.Y != 0 {
		t.Error("Coordinates wrong!", coord)
	}
}

func TestPaddelUpdate(t *testing.T) {
	paddel := getPaddelForTest(t)

	paddel.MoveDown(1)

	coord, _ := paddel.GetCoord()
	if coord.Y != 1 {
		t.Error("Coordinate after moving down is wrong:", coord.Y)
	}

	paddel.MoveUp(1)

	coord, _ = paddel.GetCoord()
	if coord.Y != 0 {
		t.Error("coordinate after moving up again is wrong:", coord.Y)
	}

	paddel.MoveUp(1)
	coord, _ = paddel.GetCoord()
	if coord.Y != 0 {
		t.Error("Boundary checking doesn't work!")
	}
}

func TestPaddelUpdateTiming(t *testing.T) {
	paddel := getPaddelForTest(t)

	ticker := time.NewTicker(time.Nanosecond)
	currentTick := <-ticker.C
	// generate more than just one event
	paddel.MoveDown(1)
	paddel.MoveDown(1)
	paddel.MoveDown(1)
	// need to get the latestmost ones
	coord, _ := paddel.GetCoord()
	oldTick := currentTick
	currentTick = <-ticker.C
	duration := currentTick.Sub(oldTick)
	if duration.Nanoseconds() > 45000000 {
		t.Fatal("duration:", duration.Nanoseconds())
	}
	if coord.Y != 3 {
		t.Fatal("Coord are off in timing test:", coord.Y)
	}

	currentTick = <-ticker.C
	coord, _ = paddel.GetCoord()
	if coord.Y != 3 {
		t.Fatal("Coord are off in second try:", coord.Y)
	}
	oldTick = currentTick
	currentTick = <-ticker.C
	duration = currentTick.Sub(oldTick)
	if duration.Nanoseconds() > 45000000 {
		t.Fatal("duration in second update:", duration.Nanoseconds())
	}
	if coord.Y != 3 {
		t.Fatal("Coord are off in second update timing test:", coord.Y)
	}
}
