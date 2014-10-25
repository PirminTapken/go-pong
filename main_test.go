package main

import (
	"testing"
)

func TestSdlError(t *testing.T) {
	err := NewSdlError("wurst")
	if err.Error() != "wurst" {
		t.Errorf("%v != %v", err.Error(), "wurst")
	}
}
