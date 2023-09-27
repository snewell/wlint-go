package filter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	errNoWordLists error = fmt.Errorf("no word lists provided")
)

func readWords(reader io.Reader, wordFn func(string) error) error {
	return wlint.Linify(reader, func(line string, count int) error {
		word := strings.TrimSpace(line)
		if len(word) > 0 && word[0] != '#' {
			return wordFn(word)
		}
		return nil
	})
}

func loadWordList(filename string, wordFn func(string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return readWords(f, wordFn)
}

func makeWordLists(lists []string, baseDir string) []string {
	ret := make([]string, len(lists))
	for index := range lists {
		ret[index] = path.Join(baseDir, lists[index])
	}
	return ret
}

func buildWordFilesList(cliFiles []string, configs []wlint.ConfigInfo[config]) ([]string, error) {
	// files provided via cli always take precendence
	if len(cliFiles) != 0 {
		return cliFiles, nil
	}

	// if nothing was provided via cli, check configs
	for index := range configs {
		if len(configs[index].Config.WordFilterConfig.WordListFiles) != 0 {
			// there are word lists defined, so build the full path
			return makeWordLists(configs[index].Config.WordFilterConfig.WordListFiles, configs[index].Dir), nil
		}
	}
	// no word lists available, error out
	return nil, errNoWordLists
}
