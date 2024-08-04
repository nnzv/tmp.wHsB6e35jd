package main

import "testing"

func TestReturn(t *testing.T) {
	if !(F() == "test") {
		t.FailNow()
	}
}
