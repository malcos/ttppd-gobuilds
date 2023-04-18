package a

import "testing"

func TestMessage(t *testing.T) {
	got := MessageFromA()
	expected := "I come from planet A"
	if got != "I come from planet A" {
		t.Errorf("Expected %v but got %v", expected, got)
	}
}
