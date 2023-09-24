package wc

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snewell/wlint-go/cmd"
	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	caseSensitive bool

	wordCountCmd = &cobra.Command{
		Use:   "word-count",
		Short: "Count words in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return countWordsCmd(args)
		},
	}
)

func identityString(s string) string {
	return s
}

func countWords(r io.Reader, mapFn func(string) string) (map[string]int, int, error) {
	ret := map[string]int{}
	totalCount := 0
	err := wlint.Wordify(r, func(word string) error {
		ret[mapFn(word)]++
		totalCount++
		return nil
	})
	if err != nil {
		return nil, 0, err
	} else {
		return ret, totalCount, nil
	}
}

func combineCounts(lhs *map[string]int, rhs *map[string]int) {
	for word, count := range *rhs {
		(*lhs)[word] += count
	}
}

func countWordsCmd(args []string) error {
	// start with the identity string
	mapFn := identityString
	if caseSensitive {
		mapFn = strings.ToLower
	}

	totalWords := map[string]int{}
	totalCount := 0
	err := wlint.FilesOrStdin(args, func(r io.Reader) error {
		counts, localCount, err := countWords(r, mapFn)
		if err != nil {
			return err
		}
		combineCounts(&totalWords, &counts)
		totalCount += localCount
		return nil
	})

	keys := []string{}
	for word := range totalWords {
		keys = append(keys, word)
	}
	sort.Strings(keys)
	for _, word := range keys {
		fmt.Printf("%v\t%v\n", word, totalWords[word])
	}
	return err
}

func init() {
	wordCountCmd.PersistentFlags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "Treat words as case sensitive")
	cmd.AddCommand(wordCountCmd)
}
