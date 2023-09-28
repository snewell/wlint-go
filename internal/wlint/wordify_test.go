package wlint

import (
	"bytes"
	"fmt"
	"testing"
)

type wordifyResult struct {
	word   string
	line   Line
	column Column
}

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

func validateWordifyResults(t *testing.T, actual []wordifyResult, expected []wordifyResult) {
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
	err := Wordify(reader, func(s string, line Line, column Column) error {
		if line != 1 {
			t.Errorf("Unexpected line (%v vs 1)", line)
		}
		if column != 1 {
			t.Errorf("Unexpected column (%v vs 1)", column)
		}
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
	err := Wordify(reader, func(s string, line Line, column Column) error {
		if line != 1 {
			t.Errorf("Unexpected line (%v vs 1)", line)
		}
		if column != 7 {
			t.Errorf("Unexpected column (%v vs 7)", column)
		}
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
	err := Wordify(reader, func(s string, line Line, column Column) error {
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
	expected := []wordifyResult{
		{
			word:   "foo",
			line:   1,
			column: 1,
		}, {
			word:   "foo",
			line:   2,
			column: 1,
		}, {
			word:   "foo",
			line:   3,
			column: 1,
		}, {
			word:   "foo",
			line:   6,
			column: 1,
		},
	}

	words := []wordifyResult{}
	err := Wordify(reader, func(s string, line Line, column Column) error {
		words = append(words, wordifyResult{
			word:   s,
			line:   line,
			column: column,
		})
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordifyResults(t, words, expected)
}

func TestWordifyMultiple(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo foo foo"))
	expected := []wordifyResult{
		{
			word:   "foo",
			line:   1,
			column: 1,
		}, {
			word:   "foo",
			line:   1,
			column: 5,
		}, {
			word:   "foo",
			line:   1,
			column: 9,
		},
	}

	words := []wordifyResult{}
	err := Wordify(reader, func(s string, line Line, column Column) error {
		words = append(words, wordifyResult{
			word:   s,
			line:   line,
			column: column,
		})
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	validateWordifyResults(t, words, expected)
}

func TestWordifyErr(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo"))
	errSomeError := fmt.Errorf("some error")
	err := Wordify(reader, func(s string, line Line, column Column) error {
		return errSomeError
	})
	if err != errSomeError {
		t.Errorf("Unexpected error: %v", err)
	}
}
