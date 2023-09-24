package filter

import (
	"testing"
)

func compareHits(t *testing.T, actual []patternMatch, expected []patternMatch) {
	if len(actual) != len(expected) {
		t.Errorf("Length mismatch (%v vs %v)", len(actual), len(expected))
	}
	maxIndex := len(actual)
	if maxIndex > len(expected) {
		maxIndex = len(expected)
	}
	for index := range actual[:maxIndex] {
		if actual[index] != expected[index] {
			t.Errorf("Hit mismatch (%v vs %v)", actual[index], expected[index])
		}
	}
}

func TestCaseSensitiveRegexHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "hello world")
	expected := []patternMatch{
		{
			match: "hello",
			index: 0,
		},
	}
	compareHits(t, hits, expected)
}

func TestInnerHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("ell", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "hello")
	expected := []patternMatch{}
	compareHits(t, hits, expected)
}

func TestCaseSensitiveRegexMultiHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "hello hello world")
	expected := []patternMatch{
		{
			match: "hello",
			index: 0,
		},
		{
			match: "hello",
			index: 6,
		},
	}
	compareHits(t, hits, expected)
}

func TestCaseSensitiveRegexMiss(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "Hello World")
	expected := []patternMatch{}
	compareHits(t, hits, expected)
}

func TestCaseInsensitiveRegexHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "Hello World")
	expected := []patternMatch{
		{
			match: "Hello",
			index: 0,
		},
	}
	compareHits(t, hits, expected)
}

func TestCaseInsensitiveRegexMultiHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := getRegexHits(r, "hElLo HeLlO world")
	expected := []patternMatch{
		{
			match: "hElLo",
			index: 0,
		},
		{
			match: "HeLlO",
			index: 6,
		},
	}
	compareHits(t, hits, expected)
}
