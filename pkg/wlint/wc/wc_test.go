package wc

import (
	"bytes"
	"strings"
	"testing"
)

func validateTotalCount(t *testing.T, counts map[string]int, expectedTotal int) {
	totalCount := 0
	for word, count := range counts {
		totalCount += count
		if count < 1 {
			t.Errorf("%v has an invalid count (%v)", word, count)
		}
	}
	if totalCount != expectedTotal {
		t.Errorf("Counts don't match expected total (%v vs %v)", totalCount, expectedTotal)
	}
}

func validateExpectedCounts(t *testing.T, actual map[string]int, expected map[string]int) {
	if len(actual) != len(expected) {
		t.Errorf("Size mismatch (%v vs %v)", len(actual), len(expected))
	}
	for word, expectedCount := range expected {
		if actualCount, found := actual[word]; found {
			if expectedCount != actualCount {
				t.Errorf("Count mismatch for %v (%v vs %v)", word, actualCount, expectedCount)
			}
		} else {
			t.Errorf("Missing word %v", word)
		}
	}
}

func TestCountEmpty(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte{})
	counts, total, err := countWords(reader, identityString)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	if len(counts) != 0 {
		t.Errorf("Word counts isn't empty")
	}
	validateTotalCount(t, counts, total)
}

func TestCountSensitive(t *testing.T) {
	reader := bytes.NewReader([]byte("foo FOO"))
	counts, total, err := countWords(reader, identityString)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	expected := map[string]int{
		"foo": 1,
		"FOO": 1,
	}
	validateExpectedCounts(t, counts, expected)
	validateTotalCount(t, counts, total)
}

func TestCountInsensitive(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo FOO"))
	counts, total, err := countWords(reader, strings.ToLower)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	expected := map[string]int{
		"foo": 2,
	}
	validateExpectedCounts(t, counts, expected)
	validateTotalCount(t, counts, total)
}
