package wlint

import (
	"bytes"
	"fmt"
	"testing"
)

func validateWordLists(t *testing.T, actual []string, expected []string) {
	if len(actual) != len(expected) {
		t.Errorf("Word slices different sizes (%v vs %v)", len(actual), len(expected))
	}
	maxIndex := len(actual)
	if maxIndex > len(expected) {
		maxIndex = len(expected)
	}
	for index := range actual[0:maxIndex] {
		if actual[index] != expected[index] {
			t.Errorf("Mismatch in actual vs expected at index %v (%v vs %v)", index, actual[index], expected[index])
		}
	}
}

func TestWordifySingle(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo"))
	expected := []string{"foo"}
	words := []string{}
	err := Wordify(reader, func(s string) error {
		words = append(words, s)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordLists(t, words, expected)
}

func TestWordifyPrefixSpace(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("  \t   foo"))
	expected := []string{"foo"}
	words := []string{}
	err := Wordify(reader, func(s string) error {
		words = append(words, s)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordLists(t, words, expected)
}

func TestWordifySuffixSpace(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo  \t   "))
	expected := []string{"foo"}
	words := []string{}
	err := Wordify(reader, func(s string) error {
		words = append(words, s)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordLists(t, words, expected)
}

func TestWordifyNewlines(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo\nfoo\nfoo\n\n\nfoo"))
	expected := []string{"foo", "foo", "foo", "foo"}
	words := []string{}
	err := Wordify(reader, func(s string) error {
		words = append(words, s)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordLists(t, words, expected)
}

func TestWordifyMultiple(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo foo foo"))
	expected := []string{"foo", "foo", "foo"}
	words := []string{}
	err := Wordify(reader, func(s string) error {
		words = append(words, s)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordLists(t, words, expected)
}

func TestWordifyErr(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo"))
	errSomeError := fmt.Errorf("some error")
	err := Wordify(reader, func(s string) error {
		return errSomeError
	})
	if err != errSomeError {
		t.Errorf("Unexpected error: %v", err)
	}
}
