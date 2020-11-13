package main

import (
	"os"

	"github.com/1efty/semtag/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
