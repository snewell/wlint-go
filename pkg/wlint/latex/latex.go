package latex

import (
	"io"
	"strings"

	"github.com/snewell/wlint-go/internal/wlint"
)

type LatexPurifier struct {
	R io.Reader
}

func (lp LatexPurifier) Wordify(wordFn func(string, wlint.Line, wlint.Column) error) error {
	return lp.Linify(func(text string, line wlint.Line) error {
		return wlint.WordifyString(text, func(word string, column wlint.Column) error {
			return wordFn(word, line, column)
		})
	})
}

func (lp LatexPurifier) Linify(lineFn func(string, wlint.Line) error) error {
	return wlint.Linify(lp.R, func(line string, count wlint.Line) error {
		uncommented := stripComments(line)
		return lineFn(uncommented, count)
	})
}

func makeLatexPurifier(r io.Reader) (wlint.Purifier, error) {
	return LatexPurifier{
		R: r,
	}, nil
}

func stripComments(line string) string {
	index := strings.Index(line, "%")
	if index != -1 {
		// found a comment
		return line[:index]
	}
	return line
}

func init() {
	err := wlint.AddPurifier("latex", makeLatexPurifier)
	if err != nil {
		panic(err)
	}
}
