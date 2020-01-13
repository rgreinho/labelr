package cmd

import (
	"log"
	"os"

	"github.com/rgreinho/labelr/labelr"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies GitHub labels",
	Long:  `Applies GitHub labels to a repository or to an organization.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the owner and repository.
		// 1. from repo
		// 2. from env var
		// 3. from CLI
		owner, repository := labelr.GetInfo()
		if owner == "" {
			o, err := cmd.Flags().GetString("owner")
			if err != nil {
				log.Fatalf("No owner name specified: %s", err)
			}
			owner = o
		}

		// Get repository.
		if repository == "" {
			r, err := cmd.Flags().GetString("repository")
			if err != nil {
				log.Fatalf("No repository name specified: %s", err)
			}
			repository = r
		}

		// Get token.
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			t, err := cmd.Flags().GetString("token")
			if err != nil {
				log.Fatalf("No token specified: %s", err)
			}
			token = t
		}

		// Get organization.
		org, _ := cmd.Flags().GetBool("org")
		if org && owner == "" {
			if owner = os.Getenv("GITHUB_ORGANIZATION"); owner == "" {
				log.Fatalf("An owner is required. Run this command from a repository, set '$GITHUB_ORGANIZATION', or use '--owner'")
			}
		}

		// Prepare client.
		g := labelr.NewLabelr(owner, repository, token)

		// Get sync flag (delete labels if need be).
		sync, _ := cmd.Flags().GetBool("sync")

		// Set label file.
		labelFile := "labels.yml"
		if len(args) == 1 {
			labelFile = args[0]
		}

		// Apply labels.
		var err error
		if org {
			err = g.ApplyToOrg(sync, labelFile, owner)
		} else {
			err = g.Apply(sync, labelFile)
		}
		if err != nil {
			log.Fatalf("Cannot apply labels: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.PersistentFlags().Bool("org", false, "Apply labels to a GitHub organization")
	applyCmd.Flags().Bool("sync", false, "Synchronize the labels")
}
