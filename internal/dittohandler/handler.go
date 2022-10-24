package dittohandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/juanvillacortac/ditto/pkg/program"
)

type Handler struct {
	programConfig  program.ProgramConfig
	generateConfig program.GenerateConfig
	schema         string
}

type client interface {
	// Read a file from a git repository
	ReadFile(owner, repo, branch, path string) (io.ReadCloser, error)
}

// From (url.(http.Request).URL).Path, extract the required information for fetching the
// schema file.
func parseSchemaPath(url string) struct{ owner, repo, branch, path string } {
	// Split the path into its components
	pathComponents := strings.Split(url, "/")

	return struct {
		owner, repo, branch, path string
	}{
		owner:  pathComponents[1],
		repo:   pathComponents[2],
		branch: pathComponents[3],
		path:   strings.Join(pathComponents[4:], "/"),
	}
}

func fetchGitHubSchema(path string) (*program.Schema, error) {
	// Parse the path into its components
	parsedPath := parseGitHubPath(path)

	// Fetch the schema file
	schema, err := program.FetchSchema(parsedPath.owner, parsedPath.repo, parsedPath.branch, parsedPath.path)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func Handle(w http.ResponseWriter, r *http.Request) error {
	// Parse the path to the schema file on GitHub
	schemaPath := parseGitHubURL(r.URL.Path)

	programConfig, generateConfig, err := parseConfigFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	program, err := program.New(programConfig, generateConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	program.Generate()
}

func fromReq(r *http.Request) (*program.ProgramConfig, *program.GenerateConfig, *program.Schema, error) {
	if r.Method != http.MethodPost {
		return nil, nil, fmt.Errorf("invalid request method %s", r.Method)
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return nil, nil, fmt.Errorf("invalid content type %s", r.Header.Get("Content-Type"))
	}

	if r.Body != nil {
		defer r.Body.Close()

		var config struct {
			ProgramConfig  *program.ProgramConfig  `json:"program_config"`
			GenerateConfig *program.GenerateConfig `json:"generate_config"`
		}

		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid request body: %s", err)
		}

		return config.ProgramConfig, config.GenerateConfig, nil
	}

	var programConfig program.ProgramConfig
	if r.URL.Query().Get("program_config") != "" {
		err := json.Unmarshal([]byte(r.URL.Query().Get("program_config")), &programConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid program config: %s", err)
		}
	}

	var generateConfig program.GenerateConfig
	if r.URL.Query().Get("generate_config") != "" {
		err := json.Unmarshal([]byte(r.URL.Query().Get("generate_config")), &generateConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid generate config: %s", err)
		}
	}

	return programConfig, generateConfig, nil
}
