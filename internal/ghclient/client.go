package ghclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Implement the dittoclient.Client interface.
type Client struct {
	owner, repo, branch, path string
}

// Parse the path to the schema file on GitHub.
func NewClient(url string) *Client {
	// Split the path into its components.
	pathComponents := strings.Split(url, "/")

	return &Client{
		owner:  pathComponents[1],
		repo:   pathComponents[2],
		branch: pathComponents[3],
		path:   strings.Join(pathComponents[4:], "/"),
	}
}

// Read a file from a git repository.
func (c *Client) ReadFile() (io.ReadCloser, error) {
	// Fetch the file.
	res, err := http.Get(c.url())
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (c *Client) url() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", c.owner, c.repo, c.branch, c.path)
}
