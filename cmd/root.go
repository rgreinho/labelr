package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is used by the build system.
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "labelr",
	Short:   "Manage your GitHub labels efficiently",
	Version: Version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("owner", "o", "", "GitHub owner name")
	rootCmd.PersistentFlags().StringP("repository", "r", "", "GitHub repository name")
	rootCmd.PersistentFlags().StringP("token", "t", "", "GitHub token")
}
