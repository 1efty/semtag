package cmd

import (
	"fmt"
	"os"

	"github.com/1efty/semtag/lib"
	"github.com/1efty/semtag/pkg/version"
	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = ".semtag.yaml"
)

var (
	repository     *git.Repository
	status         *git.Status
	firstVersion   *version.Version
	lastVersion    *version.Version
	currentVersion *version.Version
	finalVersion   *version.Version
	tags           []*version.Version

	varInitFncs []func()
	cmdInitFncs []func()
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// initConfig reads in config file
func initConfig() {
	if CfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		viper.SetConfigName(".semtag")
		viper.SetConfigType("yaml")

		// search for config in $HOME and current directory
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config file not found; nothing to do here
		} else {
			// config file was found but another error was produced; same
		}
	}
}

func initGit() {
	// get repository
	repository = lib.GetRepository()

	// retrieve all tags as lib.Version
	tags = lib.GetTagsAsVersion(repository)

	// determine first, last, current, and final version
	switch numOfTags := len(tags); numOfTags {
	case 0:
		firstVersion = &version.Version{LeadingV: false, Semver: semver.New("0.0.0")}
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
		finalVersion = lib.GetFinalVersion(repository)
	}
}

func tagAction(repository *git.Repository, tag string, dryrun bool) {
	// get status of worktree
	// exit if --force is not set and worktree contains changes
	if status, err := lib.GetStatus(repository); len(status) > 0 && !Force {
		lib.CheckIfError(err)
		lib.Info("\nThe following changes were found in the worktree:\n\n" +
			fmt.Sprintln(status) +
			"--force was not declared. Tag was not created.\n")
		os.Exit(1)
	}

	// override the tag to be created if -v flag is set
	if Version != "" {
		tag = Version
	}

	if dryrun {
		lib.Info(fmt.Sprintf("To be tagged: %s", tag))
	} else {
		lib.CreateTag(repository, tag)
	}
}

// GetConfigFile returns the current config file
func GetConfigFile() string {
	return CfgFile
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := Main(); err != nil {
		er(err)
	}
}

// RegisterCommandVar is used to register with semtag the initialization function
// for the command variable.
// Something must be returned to use the `var _ = ` trick.
func RegisterCommandVar(c func()) bool {
	varInitFncs = append(varInitFncs, c)

	return true
}

// RegisterCommandInit is used to register with px the initialization function
// for the command flags.
// Something must be returned to use the `var _ = ` trick.
func RegisterCommandInit(c func()) bool {
	cmdInitFncs = append(cmdInitFncs, c)
	return true
}

// Main starts the semtag cli
// Stupid simple initialization
func Main() error {
	// Setup all variables.
	// Setting up all the variables first will allow semtag
	// to initialize the init functions in any order
	for _, v := range varInitFncs {
		v()
	}

	// Call all plugin inits
	for _, f := range cmdInitFncs {
		f()
	}

	// Execute semtag
	return rootCmd.Execute()
}
