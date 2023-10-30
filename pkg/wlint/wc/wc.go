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

type countInfo struct {
	word  string
	count int
}

var (
	caseSensitive bool
	relativeUsage bool
	sortMethod    string

	sortMethods = map[string]func(countInfo, countInfo) bool{
		"alpha": func(lhs countInfo, rhs countInfo) bool {
			return lhs.word < rhs.word
		},
		"count": func(lhs countInfo, rhs countInfo) bool {
			return lhs.count < rhs.count
		},
	}

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

	sorter, found := sortMethods[sortMethod]
	if !found {
		return fmt.Errorf("unknown sort method: sortMehod")
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

	totalWords := 0
	ci := make([]countInfo, len(counter.counts))
	index := 0
	for word, count := range counter.counts {
		ci[index].word = word
		ci[index].count = count
		totalWords += count
		index++
	}

	sort.Slice(ci, func(lhs int, rhs int) bool {
		return sorter(ci[lhs], ci[rhs])
	})
	printer := func(index int) {
		fmt.Printf("%v\t%v\n", ci[index].word, ci[index].count)
	}
	if relativeUsage {
		totalWordsF := float64(totalWords)
		printer = func(index int) {
			fmt.Printf("%v\t%v\t%.3f%%\n", ci[index].word, ci[index].count, float64(ci[index].count)/totalWordsF*100)
		}
	}
	for index := range ci {
		printer(index)
	}
	fmt.Printf("Total words: %v\n", totalWords)
	return err
}

func init() {
	cmd.AddCommonFlags(wordCountCmd)
	wordCountCmd.PersistentFlags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "Treat words as case sensitive")
	wordCountCmd.PersistentFlags().BoolVarP(&relativeUsage, "relative", "r", false, "Include relative usage of each word")
	wordCountCmd.PersistentFlags().StringVarP(&sortMethod, "sort-method", "m", "alpha", "Method to order output (options are alpha and count)")
	cmd.AddCommand(wordCountCmd)
}
