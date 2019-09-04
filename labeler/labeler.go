package labeler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var remoteRegex = regexp.MustCompile(`(?m)github.com(?:[:/])(?P<owner>[^\/]*)/(?P<repo>[^\/]*)\.git`)

// var remoteRegex = regexp.MustCompile(`(?im)github.com(?:[:/])([^\/]*)/([^\/]*)\.git`)

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
func (g *GHClient) Apply(sync bool) error {
	// Read label file.
	labelfile := "/Users/remy/projects/aura-atx/.github/labels.yml"
	newLabels, err := ParseFile(labelfile)
	if err != nil {
		return fmt.Errorf("cannot parse label file: %s", err)
	}

	// Delete existing labels if need be.
	if sync {
		err := g.DeleteLabels()
		if err != nil {
			return fmt.Errorf("cannot delete labels: %s", err)
		}
	}

	// Go through the labels from the file.
	for _, l := range newLabels.Labels {
		color := strings.Replace(l.Color, "#", "", -1)

		// Convert labels to `github.Label`.
		ghLabel := &github.Label{
			Name:        &l.Name,
			Color:       &color,
			Description: &l.Description,
		}

		// Create labels.
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
func (g *GHClient) List() ([]*github.Label, error) {
	ctx := context.Background()
	labels, _, err := g.Client.Issues.ListLabels(ctx, g.Owner, g.Repository, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return labels, nil
}

// DeleteLabels delete all labels in a repository.
func (g *GHClient) DeleteLabels() error {
	ctx := context.Background()
	labels, err := g.List()
	if err != nil {
		return err
	}
	for _, label := range labels {
		_, err := g.Client.Issues.DeleteLabel(ctx, g.Owner, g.Repository, *label.Name)
		if err != nil {
			return fmt.Errorf("cannot delete label %q: %s", *label.Name, err)
		}
	}
	return nil
}

// ApplyToOrg applies labels to a full organization.
func (g *GHClient) ApplyToOrg(organization string, sync bool) error {
	ctx := context.Background()
	repositories := []*github.Repository{}
	if organization != "" {
		repos, _, err := g.Client.Repositories.ListByOrg(ctx, organization, &github.RepositoryListByOrgOptions{})
		if err != nil {
			return fmt.Errorf("Cannot list repositories of %q organization: %s", organization, err)
		}
		repositories = repos
	}
	for _, r := range repositories {
		g.Repository = *r.Name
		if err := g.Apply(sync); err != nil {
			return fmt.Errorf("cannot apply labels to %q organization: %s", organization, err)
		}
	}
	return nil
}

// GetInfo gets owner/repo information.
func GetInfo() (owner, repo string) {
	// Try to get the info from the repo itself.
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err == nil {
		matches := remoteRegex.FindStringSubmatch(string(out))
		if matches != nil && len(matches) == 3 {
			return matches[1], matches[2]
		}
	}

	// Try to get the info from the environment variables.
	return os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_REPOSITORY")
}
