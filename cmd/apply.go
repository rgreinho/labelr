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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		organization := os.Getenv("GITHUB_ORGANIZATION")
		if organization == "" {
			org, err := cmd.Flags().GetString("org")
			if err != nil {
				log.Fatalf("No organization name specified: %s", err)
			}
			organization = org
		}

		// Prepare client.
		g := labeler.NewGHClient(owner, repository, token)

		// Delete labels if need be.
		sync, err := cmd.Flags().GetBool("sync")

		// Apply labels.
		if organization == "" {
			err = g.Apply(sync)
		} else {
			g.Owner = organization
			err = g.ApplyToOrg(organization, sync)
		}
		if err != nil {
			log.Fatalf("Cannot apply labels: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.PersistentFlags().StringP("org", "", "", "GitHub organization")
	applyCmd.Flags().Bool("sync", false, "Synchronize the labels")
}
