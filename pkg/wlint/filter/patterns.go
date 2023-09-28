package filter

import (
	"fmt"
	"regexp"
)

type patternsBuilder struct {
	patterns map[string]*regexp.Regexp
}

func newPatternsBuilder() patternsBuilder {
	return patternsBuilder{
		patterns: map[string]*regexp.Regexp{},
	}
}

func (pb *patternsBuilder) add(word string, caseSensitive bool) error {

	if _, found := pb.patterns[word]; !found {
		pattern, err := buildRegex(word, caseSensitive)
		if err != nil {
			return err
		}
		pb.patterns[word] = pattern
	}
	return nil
}

func buildRegex(pattern string, caseSensitive bool) (*regexp.Regexp, error) {
	if caseSensitive {
		return regexp.Compile(fmt.Sprintf(`\b(%v)\b`, pattern))
	}
	return regexp.Compile(fmt.Sprintf(`(?i)\b(%v)\b`, pattern))
}

type patternMatch struct {
	match string
	index int
}

func getRegexHits(pattern *regexp.Regexp, text string, matchFn func(patternMatch) error) error {
	indexes := pattern.FindAllStringIndex(text, -1)
	if len(indexes) > 0 {
		for index := range indexes {
			lhs := indexes[index][0]
			rhs := indexes[index][1]
			err := matchFn(patternMatch{
				match: text[lhs:rhs],
				index: indexes[index][0],
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
