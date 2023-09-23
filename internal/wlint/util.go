package wlint

import (
	"io"
	"os"
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
