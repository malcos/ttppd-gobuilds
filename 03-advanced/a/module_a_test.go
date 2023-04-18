package a

import "testing"

func TestMessage(t *testing.T) {
	got := MessageFromA()
	expected := "I come from planet A"
	if got != expected {
		t.Errorf("Expected %v but got %v", expected, got)
	}
}
