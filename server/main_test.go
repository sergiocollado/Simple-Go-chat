package main

//  reference: https://go.dev/doc/tutorial/add-a-test

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
func TestReadFirstWord2(t *testing.T) {
	message := "\tHello World!"
	value := readFirstWord(message)
	want := "Hello"
	if want != value {
		t.Errorf("readFirstWord(\"Hello World!\") = \"%s\" , want match for \"%s\"", value, want)
	}
}
func TestReadFirstWord3(t *testing.T) {
	message := "  Hello World!"
	value := readFirstWord(message)
	want := "Hello"
	if want != value {
		t.Errorf("readFirstWord(\"Hello World!\") = \"%s\" , want match for \"%s\"", value, want)
	}
}
