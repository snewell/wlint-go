package filter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/snewell/wlint-go/internal/wlint"
)

func checkWordList(t *testing.T, actual []string, expected []string) {
	if len(actual) != len(expected) {
		t.Errorf("Length mismatch (%v vs %v)", len(actual), len(expected))
	}
	maxIndex := len(actual)
	if maxIndex > len(expected) {
		maxIndex = len(expected)
	}
	for index := range actual[:maxIndex] {
		if actual[index] != expected[index] {
			t.Errorf("Mismatch at index %v (%v vs %v)", index, actual[index], expected[index])
		}
	}
}

func TestReadWords(t *testing.T) {
	t.Parallel()

	words := []string{
		"hello",
		"world",
	}
	reader := bytes.NewReader([]byte(strings.Join(words, "\n")))
	readWords, err := readWords(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	checkWordList(t, readWords, words)
}

func TestReadExtraSpaces(t *testing.T) {
	t.Parallel()

	words := []string{
		"  hello",
		"world  ",
	}
	reader := bytes.NewReader([]byte(strings.Join(words, "\n")))
	readWords, err := readWords(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{
		"hello",
		"world",
	}
	checkWordList(t, readWords, expected)
}

func TestReadEmptyLines(t *testing.T) {
	t.Parallel()

	words := []string{
		"hello",
		"",
		"world",
	}
	reader := bytes.NewReader([]byte(strings.Join(words, "\n")))
	readWords, err := readWords(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{
		"hello",
		"world",
	}
	checkWordList(t, readWords, expected)
}

func TestReadCommentLines(t *testing.T) {
	t.Parallel()

	words := []string{
		"hello",
		"# a comment",
		"world",
	}
	reader := bytes.NewReader([]byte(strings.Join(words, "\n")))
	readWords, err := readWords(reader)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{
		"hello",
		"world",
	}
	checkWordList(t, readWords, expected)
}

func TestCliWordFiles(t *testing.T) {
	t.Parallel()

	cliFiles := []string{
		"foo",
		"bar",
	}
	configs := []wlint.ConfigInfo[config]{
		{
			Config: config{
				WordFilterConfig: wordFilterConfig{
					WordListFiles: []string{
						"biz",
						"buz",
					},
				},
			},
		},
	}
	wordFiles, err := buildWordFilesList(cliFiles, configs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	checkWordList(t, wordFiles, cliFiles)
}

func TestFirstConfigWordFiles(t *testing.T) {
	t.Parallel()

	cliFiles := []string{}
	configFiles := [][]string{
		{
			"biz",
			"buz",
		},
		{
			"abc",
			"xyz",
		},
	}
	configs := []wlint.ConfigInfo[config]{
		{
			Config: config{
				WordFilterConfig: wordFilterConfig{
					WordListFiles: configFiles[0],
				},
			},
			Dir: "/some/dir",
		},
		{
			Config: config{
				WordFilterConfig: wordFilterConfig{
					WordListFiles: configFiles[1],
				},
			},
		},
	}
	wordFiles, err := buildWordFilesList(cliFiles, configs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{
		"/some/dir/biz",
		"/some/dir/buz",
	}
	checkWordList(t, wordFiles, expected)
}

func TestNoWordFiles(t *testing.T) {
	t.Parallel()

	cliFiles := []string{}
	configFiles := [][]string{{}, {}}
	configs := []wlint.ConfigInfo[config]{
		{
			Config: config{
				WordFilterConfig: wordFilterConfig{
					WordListFiles: configFiles[0],
				},
			},
		},
		{
			Config: config{
				WordFilterConfig: wordFilterConfig{
					WordListFiles: configFiles[1],
				},
			},
		},
	}
	_, err := buildWordFilesList(cliFiles, configs)
	if err != errNoWordLists {
		t.Errorf("Unexpected error: %v", err)
	}
}
