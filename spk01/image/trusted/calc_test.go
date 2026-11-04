package trusted

import (
	"testing"
	"time"
)

func TestAddSuccess(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Errorf("expected 5")
	}
}

func TestInfiniteLoop(t *testing.T) {
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
