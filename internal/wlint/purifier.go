package wlint

import (
	"fmt"
	"io"
)

var (
	purifiers = map[string]MakePurifierFn{}
)

type MakePurifierFn func(io.Reader) (Purifier, error)

func AddPurifier(name string, makePurifierFn MakePurifierFn) error {
	if _, found := purifiers[name]; found {
		return fmt.Errorf("duplicate purifier named %v", name)
	}
	purifiers[name] = makePurifierFn
	return nil
}

func GetPurifier(name string) (MakePurifierFn, error) {
	if fn, found := purifiers[name]; found {
		return fn, nil
	}
	return nil, fmt.Errorf("unknown purifier %v", name)
}

type Purifier interface {
	Wordify(func(string, Line, Column) error) error
	Linify(func(string, Line) error) error
}

type NullPurifier struct {
	R io.Reader
}

func (np NullPurifier) Wordify(wordFn func(string, Line, Column) error) error {
	return Wordify(np.R, func(word string, line Line, column Column) error {
		return wordFn(word, line, column)
	})
}

func (np NullPurifier) Linify(lineFn func(string, Line) error) error {
	return Linify(np.R, func(word string, line Line) error {
		return lineFn(word, line)
	})
}

func makeNullPurifier(r io.Reader) (Purifier, error) {
	return NullPurifier{
		R: r,
	}, nil
}

func init() {
	err := AddPurifier("null", makeNullPurifier)
	if err != nil {
		panic(err)
	}
}
