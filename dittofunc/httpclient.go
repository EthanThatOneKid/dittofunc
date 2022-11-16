package dittofunc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/ethanthatonekid/dittofunc/dittofunc/dittoclient"
	"github.com/ethanthatonekid/dittofunc/dittofunc/githubclient"
	"github.com/pkg/errors"
)

type HTTPClient struct {
	*http.Client
	Origin string
}

// NewHTTPClient creates a new Client instance.
func NewHTTPClient(httpClient http.Client, origin string) *HTTPClient {
	return &HTTPClient{Client: &httpClient, Origin: origin}
}

// GenQuery is a query required by the Gen function.
type GenQuery dittoclient.GenQuery

// NewGenQuery is a query for the Gen function.
func NewGenQuery(token, owner, repo, path, ref string) *GenQuery {
	return &GenQuery{
		Token: token,
		RawFileQuery: githubclient.RawFileQuery{
			Owner: owner,
			Repo:  repo,
			Path:  path,
			Ref:   ref,
		},
	}
}

// Gen generates a new program using the given query and passing it to the given HTTP client.
func (c *HTTPClient) Gen(q *GenQuery) (*dittoclient.Output, error) {
	// Make the request URL.
	u, err := c.makeRequestURL(*q)
	if err != nil {
		return nil, err
	}

	// Make the request.
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	// Set the authorization header.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", q.Token))

	// Execute the request.
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	// Check the status code.
	if resp.StatusCode != http.StatusOK {
		if b, err := io.ReadAll(resp.Body); err == nil {
			return nil, errors.Wrapf(fmt.Errorf("unexpected response body %s", b), "unexpected status code %d", resp.StatusCode)
		}

		return nil, errors.Wrap(err, "failed to read response body")
	}

	// Decode the response.
	var output dittoclient.Output
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &output, nil
}

// makeRequestURL makes a request URL from the given path.
func (c *HTTPClient) makeRequestURL(q GenQuery) (*url.URL, error) {
	scheme, host, err := c.cutOrigin()
	if err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path.Join(q.Owner, q.Repo, q.Path),
		RawQuery: url.Values{
			"ref": []string{q.Ref},
		}.Encode(),
	}, nil
}

// cutOrigin cuts the origin into a scheme and host.
func (c *HTTPClient) cutOrigin() (string, string, error) {
	switch {
	case strings.HasPrefix(c.Origin, "http://"):
		return "http", strings.TrimPrefix(c.Origin, "http://"), nil
	case strings.HasPrefix(c.Origin, "https://"):
		return "https", strings.TrimPrefix(c.Origin, "https://"), nil
	default:
		return "", "", fmt.Errorf("invalid origin %s", c.Origin)
	}
}
