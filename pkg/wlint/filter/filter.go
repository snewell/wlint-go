package filter

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/snewell/wlint-go/cmd"
	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	caseSensitive bool
	wordLists     []string

	listFilterCmd = &cobra.Command{
		Use:   "list-filter",
		Short: "Identify filter words",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFilter(args)
		},
	}
)

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

func listFilter(args []string) error {
	totalWordList := []string{}
	for index := range wordLists {
		wordList, err := loadWordList(wordLists[index])
		if err != nil {
			return err
		}
		totalWordList = append(totalWordList, wordList...)
	}
	patternList := []*regexp.Regexp{}
	for index := range totalWordList {
		pattern, err := buildRegex(totalWordList[index], caseSensitive)
		if err != nil {
			return err
		}
		patternList = append(patternList, pattern)
	}
	err := wlint.FilesOrStdin(args, func(r io.Reader) error {
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
