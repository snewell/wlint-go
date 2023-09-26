package filter

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/snewell/wlint-go/cmd"
	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	errNoWordlists error = fmt.Errorf("no word lists provided")
	errNoWords     error = fmt.Errorf("no words in filter list")

	caseSensitive bool
	wordLists     []string

	listFilterCmd = &cobra.Command{
		Use:   "word-filter",
		Short: "Identify filter words",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFilter(args)
		},
	}
)

type wordFilterConfig struct {
	WordListFiles []string `yaml:"word_files"`
}

type config struct {
	wlint.Config     `yaml:",inline"`
	WordFilterConfig wordFilterConfig `yaml:"word_filter"`
}

func loadWordList(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ret := []string{}
	err = wlint.Linify(f, func(line string, count int) error {
		if len(line) > 0 && line[0] != '#' {
			ret = append(ret, string(line))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func makeWordLists(lists []string, baseDir string) []string {
	ret := make([]string, len(lists))
	for index := range lists {
		ret[index] = path.Join(baseDir, lists[index])
	}
	return ret
}

func listFilter(args []string) error {
	globalConfig, localConfig, err := wlint.GetAllConfigs[config]()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	// if nothing was provided via cli, check configs
	if len(wordLists) == 0 {
		if len(localConfig.Config.WordFilterConfig.WordListFiles) != 0 {
			wordLists = makeWordLists(localConfig.Config.WordFilterConfig.WordListFiles, localConfig.Dir)
		} else {
			wordLists = makeWordLists(globalConfig.Config.WordFilterConfig.WordListFiles, globalConfig.Dir)
		}
	}

	if len(wordLists) == 0 {
		return errNoWordlists
	}

	totalWordList := []string{}
	for index := range wordLists {
		wordList, err := loadWordList(wordLists[index])
		if err != nil {
			return err
		}
		totalWordList = append(totalWordList, wordList...)
	}
	if len(totalWordList) == 0 {
		return errNoWords
	}

	patternList := []*regexp.Regexp{}
	for index := range totalWordList {
		pattern, err := buildRegex(totalWordList[index], caseSensitive)
		if err != nil {
			return err
		}
		patternList = append(patternList, pattern)
	}
	err = wlint.FilesOrStdin(args, func(r io.Reader) error {
		return wlint.Linify(r, func(line string, count int) error {
			for index := range patternList {
				matches := getRegexHits(patternList[index], line)
				for matchIndex := range matches {
					fmt.Printf("%v\t%v:%v\n", matches[matchIndex].match, count, matches[matchIndex].index)
				}
			}
			return nil
		})
	})
	return err
}

func init() {
	listFilterCmd.PersistentFlags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "Treat text as case sensitive")
	listFilterCmd.PersistentFlags().StringSliceVarP(&wordLists, "word-list", "w", []string{}, "File to load word list from")
	cmd.AddCommand(listFilterCmd)
}
