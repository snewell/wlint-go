package wlint

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestEmptyLinify(t *testing.T) {
	t.Parallel()

	lines := []string{}

	errShouldntHit := fmt.Errorf("shouldn't have been hit")
	reader := bytes.NewReader([]byte(strings.Join(lines, "\n")))
	err := Linify(reader, func(string, Line) error {
		return errShouldntHit
	})
	if err != nil {
		t.Errorf("Unexpcted error: %v", err)
	}
}

func TestEmptyLines(t *testing.T) {
	t.Parallel()

	lines := []string{"", ""}

	reader := bytes.NewReader([]byte(strings.Join(lines, "\n")))
	err := Linify(reader, func(line string, count Line) error {
		if line != lines[count-1] {
			t.Errorf("Unexpected line content at index %v: %v", count-1, line)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Unexpcted error: %v", err)
	}
}

func TestNonEmptyLines(t *testing.T) {
	t.Parallel()

	lines := []string{"1", "2"}

	reader := bytes.NewReader([]byte(strings.Join(lines, "\n")))
	err := Linify(reader, func(line string, count Line) error {
		if line != lines[count-1] {
			t.Errorf("Unexpected line content at index %v: %v", count-1, line)
		}
		expected := fmt.Sprintf("%v", count)
		if line != expected {
			t.Errorf("Unexpected line content at index %v: %v", count-1, line)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Unexpcted error: %v", err)
	}
}
