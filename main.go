package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// globals
var (
	firstVersion   *Version
	lastVersion    *Version
	currentVersion *Version
	finalVersion   *Version
	tags           []*Version
	repository     *git.Repository
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	run(os.Args)
}
