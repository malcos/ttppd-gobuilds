package b

import "testing"

func TestMessage(t *testing.T) {
	got := MessageFromB()
	expected := "I come from planet B!"
	if got != expected {
		t.Errorf("Expected %v but got %v", expected, got)
	}
}
