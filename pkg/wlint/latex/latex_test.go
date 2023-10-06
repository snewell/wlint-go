package latex

import (
	"strings"
	"testing"
)

func TestStripNothing(t *testing.T) {
	input := "Hello world"
	stripped := stripComments(input)
	if stripped != input {
		t.Errorf("Unexpected result: %v (expected %v)", stripped, input)
	}
}

func TestStripComment(t *testing.T) {
	chunks := []string{
		"Hello ",
		" world",
	}
	input := strings.Join(chunks, "%")
	stripped := stripComments(input)
	if stripped != chunks[0] {
		t.Errorf("Unexpected result: %v (expected %v)", stripped, chunks[0])
	}
}

func TestStripMultipleComment(t *testing.T) {
	chunks := []string{
		"Hello ",
		" world",
		" of latex",
	}
	input := strings.Join(chunks, "%")
	stripped := stripComments(input)
	if stripped != chunks[0] {
		t.Errorf("Unexpected result: %v (expected %v)", stripped, chunks[0])
	}
}
