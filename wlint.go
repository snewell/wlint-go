package main

import (
	"github.com/snewell/wlint-go/internal/cmd"

	_ "github.com/snewell/wlint-go/pkg/wlint/latex"

	_ "github.com/snewell/wlint-go/pkg/wlint/filter"
	_ "github.com/snewell/wlint-go/pkg/wlint/wc"
)

func main() {
	cmd.Execute()
}
