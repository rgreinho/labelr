/*
Copyright © 2019 Rémy Greinhofer <remy.greinhofer@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"os"

	"github.com/rgreinho/labeler/labeler"
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
		owner, repository := labeler.GetInfo()
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
		g := labeler.NewLabeler(owner, repository, token)

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
