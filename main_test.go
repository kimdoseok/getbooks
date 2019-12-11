package main

import (
	"testing"
)

func TestProcess(t *testing.T) {
	if !Process() {
		t.Error("Process should return True all the time")
	}
}
