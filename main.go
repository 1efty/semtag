package main

import (
	"fmt"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
)

// globals
var (
	firstVersion   *semver.Version
	lastVersion    *semver.Version
	currentVersion *semver.Version
	finalVersion   *semver.Version
	tags           []*semver.Version
	repository     *git.Repository
	validScopes    = []string{"patch", "minor", "major"}
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	run(os.Args)
}
