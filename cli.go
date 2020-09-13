package main

import (
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

func setup(app *cli.App) {
	app.Name = "semtag"
	app.Usage = "Tag your repository according to Semantic Versioning."
	app.Version = "v0.0.3"
	app.Compiled = time.Now()

	// some commands share the same arguments
	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "scope",
			Aliases: []string{"s"},
			EnvVars: []string{"SEMTAG_SCOPE"},
			Value:   "patch",
			Usage: `The scope that must be increased, can be major, minor or patch.
               The resulting version will match X.Y.Z(-PRERELEASE)(+BUILD)
               where X, Y and Z are positive integers, PRERELEASE is an optional
               string composed of alphanumeric characters describing if the build is
               a release candidate, alpha or beta version, with a number.
               BUILD is also an optional string composed of alphanumeric
               characters and hyphens.
               Setting the scope as 'auto', the script will chose the scope between
               'minor' and 'patch', depending on the amount of lines added (<10% will
               choose patch).`,
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			EnvVars: []string{"SEMTAG_VERSION"},
			Usage: `Specifies manually the version to be tagged, must be a valid semantic version
			   in the format X.Y.Z where X, Y and Z are positive integers.`,
		},
		&cli.StringFlag{
			Name:    "metadata",
			Aliases: []string{"m"},
			EnvVars: []string{"SEMTAG_METADATA"},
			Usage:   "Specifies the metadata (+BUILD) for the version.",
		},
		&cli.BoolFlag{
			Name:    "output",
			Aliases: []string{"o"},
			EnvVars: []string{"SEMTAG_OUTPUT"},
			Value:   false,
			Usage:   "Output the version only, shows the bumped version, but doesn't tag.",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			EnvVars: []string{"SEMTAG_FORCE"},
			Value:   false,
			Usage:   "Forces to tag, even if there are un-staged or un-committed changes.",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "get",
			Usage:  "Returns both current and final version and last tagged version.",
			Action: getAction,
		},
		{
			Name:   "getlast",
			Usage:  "Returns the latest tagged version.",
			Action: getLastAction,
		},
		{
			Name:   "getfinal",
			Usage:  "Returns latest tagged final version.",
			Action: getFinalAction,
		},
		{
			Name: "getcurrent",
			Usage: `Returns the current version, based on the latest one, if there are un-committed or
               un-staged changes, they will be reflected in the version, adding the number of
               pending commits, current branch and commit hash.`,
			Action: getCurrentAction,
		},
		{
			Name:   "final",
			Usage:  "Tags the current ref as a final version, this only can be done on the master branch.",
			Action: finalAction,
			Flags:  commonFlags,
		},
		{
			Name:   "candidate",
			Usage:  "Tags the current ref as a release candidate, the tag will contain all the commits from the last final version.",
			Action: candidateAction,
			Flags:  commonFlags,
		},
		{
			Name:   "alpha",
			Usage:  "tags the current ref as an alpha version, the tag will contain all the commits from the last final version.",
			Action: alphaAction,
			Flags:  commonFlags,
		},
		{
			Name:   "beta",
			Usage:  "Tags the current ref as a beta version, the tag will contain all the commits from the last final version.",
			Action: betaAction,
			Flags:  commonFlags,
		},
	}

	// run this before every command
	app.Before = before
}

/// before ... run this before every CLI command
func before(c *cli.Context) error {
	// get repository
	repository = getRepository()

	// fetch default remote
	repository.Fetch(&git.FetchOptions{})

	// get and set tags
	tags = getTagsAsSemver(repository)

	// set globals
	switch numOfTags := len(tags); numOfTags {
	case 0:
		firstVersion = &Version{leadingV: false, semver: semver.New("0.0.0")}
		lastVersion = firstVersion
		currentVersion = firstVersion
		finalVersion = firstVersion
	case 1:
		firstVersion = tags[0]
		lastVersion = firstVersion
		currentVersion = firstVersion
		finalVersion = firstVersion
	default:
		firstVersion = tags[0]
		lastVersion = tags[len(tags)-1]
		currentVersion = lastVersion
		finalVersion = lastVersion
	}

	return nil
}

// run ... run CLI app
func run(args []string) error {
	app := cli.NewApp()
	setup(app)
	app.Run(args)
	return nil
}
