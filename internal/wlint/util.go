package wlint

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

const (
	rightSingleQuote string = "’"
)

var (
	wordPattern *regexp.Regexp
)

type Line int

type Column int

func FilesOrStdin(args []string, readerHandler func(io.Reader) error) error {
	if len(args) == 0 {
		return readerHandler(os.Stdin)
	} else {
		for index := range args {
			f, err := os.Open(args[index])
			if err != nil {
				return err
			}
			err = readerHandler(f)
			if err != nil {
				f.Close()
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func linifyInternal(reader *bufio.Reader, lineFunc func(string, Line) error) error {
	fullString := ""
	lineCount := Line(1)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == nil {
			fullString += string(line)
			if !isPrefix {
				err := lineFunc(fullString, lineCount)
				if err != nil {
					return err
				}
				fullString = ""
				lineCount++
			}
		} else if err == io.EOF {
			return nil
		} else {
			return err
		}
	}
}

func Linify(r io.Reader, lineFunc func(string, Line) error) error {
	return linifyInternal(bufio.NewReader(r), lineFunc)
}

func Wordify(reader io.Reader, wordFunc func(string, Line, Column) error) error {
	return Linify(reader, func(line string, count Line) error {
		return WordifyString(line, func(word string, column Column) error {
			return wordFunc(word, count, column)
		})
	})
}

func WordifyString(text string, wordFunc func(string, Column) error) error {
	indexes := wordPattern.FindAllStringIndex(string(text), -1)
	if len(indexes) != 0 {
		for index := range indexes {
			lhs := indexes[index][0]
			rhs := indexes[index][1]
			err := wordFunc(text[lhs:rhs], Column(lhs+1))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	wordPattern = regexp.MustCompile(fmt.Sprintf(`\b([\w\-\'%v]+)\b`, rightSingleQuote))
}
