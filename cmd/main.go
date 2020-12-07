package cmd

import (
	"fmt"
	"os"

	"github.com/1efty/semtag/lib"
	repo "github.com/1efty/semtag/pkg/git"
	"github.com/1efty/semtag/pkg/version"
	"github.com/coreos/go-semver/semver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = ".semtag.yaml"
)

var (
	repository *repo.Repo

	varInitFncs []func()
	cmdInitFncs []func()
)

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
	repository = repo.New(".")
}

func tagAction(repository *repo.Repo, tag string, dryrun bool) {
	// get status of worktree
	// exit if --force is not set and worktree contains changes
	if status := repository.Status; len(status) > 0 && !Force {
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
		repository.CreateTag(tag)
	}
}

func bumpVersion(v *version.Version, scope string, preRelease string, metadata string) (*version.Version, error) {
	newVersion := v
	err := newVersion.Bump(scope)
	lib.CheckIfError(err)

	// set pre-release and metadata
	newVersion.Semver.PreRelease = semver.PreRelease(preRelease)
	newVersion.Semver.Metadata = metadata

	return v, nil
}
