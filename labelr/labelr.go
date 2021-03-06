package labelr

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

// Labelr defines a GitHub client.
type Labelr struct {
	Owner      string
	Repository string
	Client     *github.Client
}

// NewLabelr creates a new GitHub client.
func NewLabelr(owner, repository, token string) *Labelr {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	g := github.NewClient(tc)
	return &Labelr{
		Owner:      owner,
		Repository: repository,
		Client:     g,
	}
}

// Apply retrieves the collaborator information.
func (l *Labelr) Apply(sync bool, labelFile string) error {
	// Read label file.
	newLabels, err := ParseFile(labelFile)
	if err != nil {
		return fmt.Errorf("cannot parse label file: %s", err)
	}

	// Delete existing labels if need be.
	if sync {
		if err := l.DeleteLabels(); err != nil {
			return fmt.Errorf("cannot delete labels: %s", err)
		}
	}

	// Go through the labels from the file.
	for _, label := range newLabels.Labels {
		color := strings.Replace(label.Color, "#", "", -1)

		// Convert labels to `github.Label`.
		ghLabel := &github.Label{
			Name:        &label.Name,
			Color:       &color,
			Description: &label.Description,
		}

		// Create labels.
		ctx := context.Background()
		_, r, err := l.Client.Issues.CreateLabel(ctx, l.Owner, l.Repository, ghLabel)

		// Update the label if the it already exists.
		if r.StatusCode == 422 {
			if _, _, err := l.Client.Issues.EditLabel(ctx, l.Owner, l.Repository, *ghLabel.Name, ghLabel); err != nil {
				return fmt.Errorf("cannot update label: %s", err)
			}
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// List lists existing labels.
func (l *Labelr) List() ([]*github.Label, error) {
	ctx := context.Background()
	labels, _, err := l.Client.Issues.ListLabels(ctx, l.Owner, l.Repository, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return labels, nil
}

// DeleteLabels delete all labels in a repository.
func (l *Labelr) DeleteLabels() error {
	ctx := context.Background()
	labels, err := l.List()
	if err != nil {
		return err
	}
	for _, label := range labels {
		if _, err := l.Client.Issues.DeleteLabel(ctx, l.Owner, l.Repository, *label.Name); err != nil {
			return fmt.Errorf("cannot delete label %q: %s", *label.Name, err)
		}
	}
	return nil
}

// ApplyToOrg applies labels to a full organization.
func (l *Labelr) ApplyToOrg(sync bool, labelFile, organization string) error {
	ctx := context.Background()
	repositories := []*github.Repository{}
	if organization != "" {
		repos, _, err := l.Client.Repositories.ListByOrg(ctx, organization, &github.RepositoryListByOrgOptions{})
		if err != nil {
			return fmt.Errorf("Cannot list repositories of %q organization: %s", organization, err)
		}
		repositories = repos
	}
	for _, r := range repositories {
		l.Repository = *r.Name
		if err := l.Apply(sync, labelFile); err != nil {
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
		if matches := remoteRegex.FindStringSubmatch(string(out)); len(matches) == 3 {
			return matches[1], matches[2]
		}
	}

	// Try to get the info from the environment variables.
	return os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_REPOSITORY")
}
