package labeler

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// GHClient defines a GitHub client.
type GHClient struct {
	Owner      string
	Repository string
	Client     *github.Client
}

// NewGHClient creates a new GitHub client.
func NewGHClient(owner, repository, token string) *GHClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	g := github.NewClient(tc)
	return &GHClient{
		Owner:      owner,
		Repository: repository,
		Client:     g,
	}
}

// Apply retrieves the collaborator information.
func (g *GHClient) Apply() error {
	// Read label file.
	newLabels, err := ParseFile("/Users/remy/projects/aura-atx/.github/labels.yml")
	if err != nil {
		return fmt.Errorf("cannot parse label file: %s", err)
	}

	// Create labels.
	for _, l := range newLabels.Labels {
		color := strings.Replace(l.Color, "#", "", -1)

		ghLabel := &github.Label{
			Name:        &l.Name,
			Color:       &color,
			Description: &l.Description,
		}

		ctx := context.Background()
		_, r, err := g.Client.Issues.CreateLabel(ctx, g.Owner, g.Repository, ghLabel)

		// Ignore error if the label already exists.
		if r.StatusCode == 422 {
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// List lists existing labels.
func (g *GHClient) List() error {
	ctx := context.Background()
	labels, _, err := g.Client.Issues.ListLabels(ctx, g.Owner, g.Repository, &github.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("%s", labels)
	return nil
}
