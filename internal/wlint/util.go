package wlint

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

const (
	rightSingleQuote string = "’"
)

var (
	wordPattern *regexp.Regexp
)

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

func linifyInternal(reader *bufio.Reader, lineFunc func(string, int) error) error {
	fullString := ""
	lineCount := 1
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

func Linify(r io.Reader, lineFunc func(string, int) error) error {
	return linifyInternal(bufio.NewReader(r), lineFunc)
}

func wordifyInternal(reader *bufio.Reader, wordFunc func(string) error) error {
	return linifyInternal(reader, func(line string, count int) error {
		results := wordPattern.FindAllStringSubmatch(string(line), -1)
		for index := range results {
			for _, word := range results[index][1:] {
				err := wordFunc(word)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func Wordify(r io.Reader, wordFunc func(string) error) error {
	return wordifyInternal(bufio.NewReader(r), wordFunc)
}

func init() {
	var err error

	wordPattern, err = regexp.Compile(fmt.Sprintf(`\b([\w\-\'%v]+)\b`, rightSingleQuote))
	if err != nil {
		log.Fatalf("Error compiling word marker regex: %v", err)
	}

}
