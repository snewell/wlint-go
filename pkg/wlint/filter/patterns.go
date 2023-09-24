package filter

import (
	"fmt"
	"regexp"
)

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

func getRegexHits(pattern *regexp.Regexp, text string) []patternMatch {
	indexes := pattern.FindAllStringIndex(text, -1)
	if len(indexes) > 0 {
		ret := make([]patternMatch, len(indexes))
		matches := pattern.FindAllStringSubmatch(text, -1)
		for index := range indexes {
			ret[index].match = matches[index][0]
			ret[index].index = indexes[index][0]
		}
		return ret
	}
	return []patternMatch{}
}