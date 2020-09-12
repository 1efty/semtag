package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func getAction(c *cli.Context) error {
	info(fmt.Sprintf("Current final version: %v", finalVersion))
	info(fmt.Sprintf("Last tagged version: %v", lastVersion))
	return nil
}

func getLastAction(c *cli.Context) error {
	info(lastVersion.String())
	return nil
}

func getFinalAction(c *cli.Context) error {
	info(finalVersion.String())
	return nil
}

func getCurrentAction(c *cli.Context) error {
	info(currentVersion.String())
	return nil
}

func finalAction(c *cli.Context) error {
	v, err := bumpVersion(lastVersion, c.String("scope"), "", "")
	checkIfError(err)

	createTag(repository, v.String())

	return nil
}

func candidateAction(c *cli.Context) error {
	v, err := bumpVersion(lastVersion, c.String("scope"), "rc", c.String("metadata"))
	checkIfError(err)

	createTag(repository, v.String())

	return nil
}

func alphaAction(c *cli.Context) error {
	v, err := bumpVersion(lastVersion, c.String("scope"), "alpha", c.String("metadata"))
	checkIfError(err)

	createTag(repository, v.String())

	return nil
}

func betaAction(c *cli.Context) error {
	v, err := bumpVersion(lastVersion, c.String("scope"), "beta", c.String("metadata"))
	checkIfError(err)

	createTag(repository, v.String())

	return nil
}
