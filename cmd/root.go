package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5"

	"github.com/1efty/semtag/lib"
)

var (
	// CfgFile is path to config file
	CfgFile string
	// Scope represents what kind of bump to perform
	Scope string
	// Force tag application
	Force bool
	// Output only
	Output bool
	// Version represents the specific tag
	Version string
	// Metadata represents suffix to append to tag
	Metadata string

	repository     *git.Repository
	status         *git.Status
	firstVersion   *lib.Version
	lastVersion    *lib.Version
	currentVersion *lib.Version
	finalVersion   *lib.Version
	tags           []*lib.Version

	rootCmd = &cobra.Command{
		Use:   "semtag",
		Short: "Tag your repository according to Semantic Versioning.",
		Long:  ``,
	}
)

// Execute the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Output, "output", "o", false, "Output the version only, shows the bumped version, but doesn't perform tag.")
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Forces tagging, even if there are un-staged or un-commited changes.")

	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file (default is \"$HOME/.semtag.yaml\")")
	rootCmd.PersistentFlags().StringVarP(&Metadata, "metadata", "m", "", "Specifies the metadata (+BUILD) for the version.")
	rootCmd.PersistentFlags().StringVarP(&Version, "version", "v", "", `Specifies manually the version to be tagged, must be a valid semantic version
				 in the format X.Y.Z where X, Y and Z are positive integers.`)
	rootCmd.PersistentFlags().StringVarP(&Scope, "scope", "s", "",
		`The scope that must be increased, can be major, minor or patch.
		The resulting version will match X.Y.Z(-PRERELEASE)(+BUILD)
		where X, Y and Z are positive integers, PRERELEASE is an optional
		string composed of alphanumeric characters describing if the build is
		a release candidate, alpha or beta version, with a number.
		BUILD is also an optional string composed of alphanumeric
		characters and hyphens.
		Setting the scope as 'auto', the script will chose the scope between
		'minor' and 'patch', depending on the amount of lines added (<10% will
		choose patch).`)

	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initGit() {
	// get repository
	repository = lib.GetRepository()

	// retrieve all tags as lib.Version
	tags = lib.GetTagsAsVersion(repository)

	// determine first, last, current, and final version
	switch numOfTags := len(tags); numOfTags {
	case 0:
		firstVersion = &lib.Version{LeadingV: false, Semver: semver.New("0.0.0")}
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

func initConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
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

	if dryrun {
		lib.Info(fmt.Sprintf("To be tagged: %s", tag))
	} else {
		lib.CreateTag(repository, tag)
	}
}
