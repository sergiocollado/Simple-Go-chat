package main

import (
	"testing"
)

func TestReadFirstWord1(t *testing.T) {
	message := "Hello World!"
	value := readFirstWord(message)
	want := "Hello"
	if want != value {
		t.Errorf("readFirstWord(\"Hello World!\") = \"%s\" , want match for \"%s\"", value, want)
	}
}
