package wc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/snewell/wlint-go/cmd"
	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	wordPattern *regexp.Regexp

	caseSensitive bool

	wordCountCmd = &cobra.Command{
		Use:   "word-count",
		Short: "Count words in a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return countWordsCmd(args)
		},
	}
)

const (
	rightSingleQuote string = "’"
)

func countWords(r io.Reader) (map[string]int, int, error) {
	reader := bufio.NewReader(r)
	ret := map[string]int{}
	totalCount := 0
	for {
		line, _, err := reader.ReadLine()
		if err == nil {
			results := wordPattern.FindAllStringSubmatch(string(line), -1)
			for index := range results {
				for _, word := range results[index][1:] {
					if !caseSensitive {
						word = strings.ToLower(word)
					}
					ret[word]++
					totalCount++
				}
			}
		} else if err == io.EOF {
			return ret, totalCount, nil
		} else {
			return nil, 0, err
		}
	}
}

func combineCounts(lhs *map[string]int, rhs *map[string]int) {
	for word, count := range *rhs {
		(*lhs)[word] += count
	}
}

func countWordsCmd(args []string) error {
	totalWords := map[string]int{}
	totalCount := 0
	err := wlint.FilesOrStdin(args, func(r io.Reader) error {
		counts, localCount, err := countWords(r)
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
	var err error

	wordPattern, err = regexp.Compile(fmt.Sprintf(`\b([\w\-\'%v]+)\b`, rightSingleQuote))
	if err != nil {
		log.Fatalf("Error compiling word count regex: %v", err)
	}

	wordCountCmd.PersistentFlags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "Treat words as case sensitive")
	cmd.AddCommand(wordCountCmd)
}
