package main

import "testing"

func TestParse(t *testing.T) {
	o, err := parse([]string{"--density", "0.4", "--speed=20", "--seed", "42", "--wrap=false", "--paused"})
	if err != nil {
		t.Fatal(err)
	}
	if o.density != .4 || o.speed != 20 || o.seed != 42 || o.wrap || !o.paused {
		t.Fatalf("parse() = %+v", o)
	}
}

func TestParseRejectsInvalidValues(t *testing.T) {
	for _, args := range [][]string{{"--density=2"}, {"--speed=0"}, {"--width=-1"}, {"extra"}, {"--intro", "--no-intro"}} {
		if _, err := parse(args); err == nil {
			t.Errorf("parse(%q) unexpectedly succeeded", args)
		}
	}
}
