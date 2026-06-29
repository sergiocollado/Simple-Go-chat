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

func TestJOIN(t *testing.T) {
	message := "   JOIN Ana"
	value := checkJOIN(message)
	want := "Ana"
	if want != value {
		t.Errorf("checkJOIN(\"   JOIN Ana\") = \"%s\" , want match for \"%s\"", value, want)
	}
}

func TestJOIN2(t *testing.T) {
	message := "   JOIN   Bob"
	value := checkJOIN(message)
	want := "Bob"
	if want != value {
		t.Errorf("checkJOIN(\"   JOIN   Bob\") = \"%s\" , want match for \"%s\"", value, want)
	}
}

func TestLEAVE(t *testing.T) {
	message := "   LEAVE "
	value := checkLEAVE(message)
	want := true
	if want != value {
		t.Errorf("checkLEAVE(\"   LEAVE\") = \"%t\" , want match for \"%t\"", value, want)
	}
}
func TestWHO(t *testing.T) {
	message := "   WHO "
	value := checkWHO(message)
	want := true
	if want != value {
		t.Errorf("checkWHO(\"   WHO\") = \"%t\" , want match for \"%t\"", value, want)
	}
}

func TestHELP(t *testing.T) {
	message := "   HELP "
	value := checkHELP(message)
	want := true
	if want != value {
		t.Errorf("checkHELP(\"   HELP\") = \"%t\" , want match for \"%t\"", value, want)
	}
}
