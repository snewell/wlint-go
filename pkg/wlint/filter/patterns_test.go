package filter

import (
	"fmt"
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
	hits := []patternMatch{}
	err = getRegexHits(r, "hello world", func(pm patternMatch) error {
		hits = append(hits, pm)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
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
	errShouldntHit := fmt.Errorf("shouldn't hit")
	err = getRegexHits(r, "hello", func(pm patternMatch) error {
		return errShouldntHit
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCaseSensitiveRegexMultiHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := []patternMatch{}
	err = getRegexHits(r, "hello hello world", func(pm patternMatch) error {
		hits = append(hits, pm)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
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
	errShouldntHit := fmt.Errorf("shouldn't be hit")
	err = getRegexHits(r, "Hello World", func(patternMatch) error {
		return errShouldntHit
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCaseInsensitiveRegexHit(t *testing.T) {
	t.Parallel()

	r, err := buildRegex("hello", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hits := []patternMatch{}
	err = getRegexHits(r, "Hello World", func(pm patternMatch) error {
		hits = append(hits, pm)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
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
	hits := []patternMatch{}
	err = getRegexHits(r, "hElLo HeLlO world", func(pm patternMatch) error {
		hits = append(hits, pm)
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
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

func TestBuildSinglePattern(t *testing.T) {
	t.Parallel()

	pb := newPatternsBuilder()
	err := pb.add("hello", true)
	if err != nil {
		t.Errorf("Unepxected error: %v", err)
	}
	if len(pb.patterns) != 1 {
		t.Errorf("Mismatch between pattern size and inputs (%v vs 1)", len(pb.patterns))
	}
}

func TestBuildDuplicatePattern(t *testing.T) {
	t.Parallel()

	words := []string{
		"hello",
		"hello",
	}
	pb := newPatternsBuilder()
	for _, word := range words {
		err := pb.add(word, true)
		if err != nil {

			t.Errorf("Unepxected error: %v", err)
		}
	}
	if len(pb.patterns) != 1 {
		t.Errorf("Unexpected size for paterns (%v vs 1)", len(words))
	}
}
