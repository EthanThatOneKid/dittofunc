package githubclient

import (
	"context"

	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

// Client wraps around the GitHub client.
type Client struct {
	*github.Client
	ctx context.Context
}

// New creates a new Client instance.
func New(ctx context.Context, token string) *Client {
	return &Client{
		Client: github.NewClient(
			oauth2.NewClient(
				ctx,
				oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: token},
				),
			),
		),
		ctx: ctx,
	}
}

// RawFileQuery contains the location of a file on GitHub and the token.
type RawFileQuery struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Path  string `json:"path"`
	Ref   string `json:"ref"`
}

// RawFile returns the raw file content of a file on GitHub.
func (c *Client) RawFile(q RawFileQuery) (string, error) {
	file, _, _, err := c.Client.Repositories.GetContents(c.ctx, q.Owner, q.Repo, q.Path, &github.RepositoryContentGetOptions{
		Ref: q.Ref,
	})
	if err != nil {
		return "", err
	}
	return file.GetContent()
}
