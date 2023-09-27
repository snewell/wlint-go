package wc

import (
	"bytes"
	"strings"
	"testing"
)

func validateTotalCount(t *testing.T, wc *wordCounter) {
	totalCount := 0
	for word, count := range wc.counts {
		totalCount += count
		if count < 1 {
			t.Errorf("%v has an invalid count (%v)", word, count)
		}
	}
	if totalCount != wc.totalCount {
		t.Errorf("Counts don't match expected total (%v vs %v)", totalCount, wc.totalCount)
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
	counter := makeWordCounter(identityString)
	err := countWords(&counter, reader)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	if len(counter.counts) != 0 {
		t.Errorf("Word counts isn't empty")
	}
	validateTotalCount(t, &counter)
}

func TestCountSensitive(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo FOO"))
	counter := makeWordCounter(identityString)
	err := countWords(&counter, reader)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	expected := map[string]int{
		"foo": 1,
		"FOO": 1,
	}
	validateExpectedCounts(t, counter.counts, expected)
	validateTotalCount(t, &counter)
}

func TestCountInsensitive(t *testing.T) {
	t.Parallel()

	reader := bytes.NewReader([]byte("foo FOO"))
	counter := makeWordCounter(strings.ToLower)
	err := countWords(&counter, reader)
	if err != nil {
		t.Errorf("Unexpected errors")
	}
	expected := map[string]int{
		"foo": 2,
	}
	validateExpectedCounts(t, counter.counts, expected)
	validateTotalCount(t, &counter)
}
