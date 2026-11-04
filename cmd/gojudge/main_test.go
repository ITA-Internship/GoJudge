package main

import "testing"

func TestGreeting(t *testing.T) {
	defGreeting := "Yo!"
	if got := greeting(); got != defGreeting {
		t.Fatalf("greeting() = %q, want = %q", got, defGreeting)
	}
}
