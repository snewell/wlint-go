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

type mapFn func(string) string

type wordCounter struct {
	counts     map[string]int
	wordMapFn  mapFn
	totalCount int
}

func makeWordCounter(fn mapFn) wordCounter {
	return wordCounter{
		counts:    map[string]int{},
		wordMapFn: fn,
	}
}

func (wc *wordCounter) add(word string) {
	actualWord := wc.wordMapFn(word)
	wc.counts[actualWord]++
	wc.totalCount++
}

func identityString(s string) string {
	return s
}

func countWords(wc *wordCounter, purifier wlint.Purifier) error {
	return purifier.Wordify(func(word string, line wlint.Line, column wlint.Column) error {
		wc.add(word)
		return nil
	})
}

func countWordsCmd(args []string) error {
	configs, err := wlint.GetAllConfigs[wlint.Config]()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	makePurifier, err := wlint.FindPurifier(cmd.Purifier, configs)
	if err != nil {
		return err
	}

	// start with the identity string
	mapFn := identityString
	if !caseSensitive {
		mapFn = strings.ToLower
	}

	counter := makeWordCounter(mapFn)
	err = wlint.FilesOrStdin(args, func(r io.Reader) error {
		purifier, err := makePurifier(r)
		if err != nil {
			return err
		}
		return countWords(&counter, purifier)
	})
	if err != nil {
		return err
	}

	keys := make([]string, len(counter.counts))
	index := 0
	for word := range counter.counts {
		keys[index] = word
		index++
	}
	sort.Strings(keys)
	for _, word := range keys {
		fmt.Printf("%v\t%v\n", word, counter.counts[word])
	}
	return err
}

func init() {
	cmd.AddCommonFlags(wordCountCmd)
	wordCountCmd.PersistentFlags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "Treat words as case sensitive")
	cmd.AddCommand(wordCountCmd)
}
