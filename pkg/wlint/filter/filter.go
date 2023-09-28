package filter

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/snewell/wlint-go/cmd"
	"github.com/snewell/wlint-go/internal/wlint"
)

var (
	errNoWords error = fmt.Errorf("no words in filter list")

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

func listFilter(args []string) error {
	configs, err := wlint.GetAllConfigs[config]()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	wordFiles, err := buildWordFilesList(wordLists, configs)
	if err != nil {
		return err
	}

	pb := newPatternsBuilder()
	for index := range wordFiles {
		err := loadWordList(wordFiles[index], func(word string) error {
			return pb.add(word, caseSensitive)
		})
		if err != nil {
			return err
		}
	}
	if len(pb.patterns) == 0 {
		return errNoWords
	}

	err = wlint.FilesOrStdin(args, func(r io.Reader) error {
		return wlint.Linify(r, func(line string, count wlint.Line) error {
			for _, pattern := range pb.patterns {
				getRegexHits(pattern, line, func(pm patternMatch) error {
					fmt.Printf("%v\t%v:%v\n", pm.match, count, pm.index)
					return nil
				})
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
