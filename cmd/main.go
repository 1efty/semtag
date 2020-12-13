package cmd

import (
	"os"

	repo "github.com/1efty/semtag/pkg/git"
	"github.com/1efty/semtag/pkg/utils"
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
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := Main()
	utils.CheckIfError(err)
}

// Main starts the semtag cli
// Stupid simple initialization
func Main() error {
	// Execute semtag
	return rootCmd.Execute()
}

// initConfig reads in config file
func initConfig() {
	if CfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// find home directory.
		home, err := homedir.Dir()
		utils.CheckIfError(err)

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
		utils.Info("\nThe following changes were found in the worktree:\n\n%s\n--force was not declared. Tag was not created.\n", status)
		os.Exit(1)
	}

	// override the tag to be created if -v flag is set
	if Version != "" {
		tag = Version
	}

	if dryrun {
		utils.Info("To be tagged: %s", tag)
	} else {
		repository.CreateTag(tag)
	}
}

func bumpVersion(v *version.Version, scope string, preRelease string, metadata string) (*version.Version, error) {
	newVersion := v
	err := newVersion.Bump(scope)
	utils.CheckIfError(err)

	// set pre-release and metadata
	newVersion.Semver.PreRelease = semver.PreRelease(preRelease)
	newVersion.Semver.Metadata = metadata

	return v, nil
}
