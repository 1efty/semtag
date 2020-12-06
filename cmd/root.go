package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
)

var rootCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	rootCmd = &cobra.Command{
		Use:   "semtag",
		Short: "Tag your repository according to Semantic Versioning",
	}
})

var _ = RegisterCommandInit(func() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Output, "output", "o", false, "Output the version only, shows the bumped version, but doesn't perform tag.")
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Forces tagging, even if there are un-staged or un-commited changes.")

	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "Specifies which config file to use")
	rootCmd.PersistentFlags().StringVarP(&Metadata, "metadata", "m", "", "Specifies the metadata (+BUILD) for the version.")
	rootCmd.PersistentFlags().StringVarP(&Version, "version", "v", "", `Specifies manually the version to be tagged, must be a valid semantic version
				 in the format X.Y.Z where X, Y and Z are positive integers.`)
	rootCmd.PersistentFlags().StringVarP(&Scope, "scope", "s", "patch",
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

	// bind flags to config
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	viper.BindPFlag("metadata", rootCmd.PersistentFlags().Lookup("metadata"))
	viper.BindPFlag("scope", rootCmd.PersistentFlags().Lookup("scope"))
})
