// dittoclient is a client for the Ditto API.
package dittoclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/ethanthatonekid/dittofunc/dittofunc/githubclient"
	"github.com/juanvillacortac/ditto/pkg/ast"
	"github.com/juanvillacortac/ditto/pkg/generators"
	"github.com/juanvillacortac/ditto/pkg/parser/proto"
	"github.com/juanvillacortac/ditto/pkg/parser/yaml"
	"github.com/juanvillacortac/ditto/pkg/program"
	"github.com/pkg/errors"
	yamlio "gopkg.in/yaml.v2"
)

// Client generates new programs.
type Client struct {
}

// NewClient creates a new Client instance.
func NewClient() *Client {
	return &Client{}
}

// GenQuery is the query for generating a new program.
type GenQuery struct {
	githubclient.RawFileQuery
	Token string
}

// Output is the output of Generate.
type Output struct {
	Files []generators.OutputFile `json:"files"`
}

// Gen generates a new program.
func (c *Client) Gen(ctx context.Context, q GenQuery) (*Output, error) {
	// Setup the Github client.
	githubClient := githubclient.NewClient(ctx, q.Token)

	// Get the program config.
	p, err := getProgram(githubClient, getProgramConfigQuery{q.RawFileQuery})
	if err != nil {
		return nil, err
	}

	// Generate the program.
	files, err := genFiles(p.Root, p.GenerateConfigs)
	if err != nil {
		return nil, err
	}

	return &Output{Files: files}, nil
}

// getConfigQuery is the query for getting the program config file.
type getProgramConfigQuery struct {
	githubclient.RawFileQuery
}

// ProgramConfig is the program config.
type Program struct {
	program.ProgramConfig
	GenerateConfigs []generators.GenerateConfig `json:"generate_configs"`
	Root            *ast.RootNode               `json:"root"`
}

// getProgram gets the required program config files.
func getProgram(githubClient *githubclient.Client, q getProgramConfigQuery) (*Program, error) {
	f, err := githubClient.RawFile(q.RawFileQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get program config file from github")
	}

	// Parse the program config file.
	p, err := parseProgramConfigFile([]byte(f), path.Ext(q.Path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse program config file")
	}

	// Get the generate configs.
	generateConfigs, err := getGenerateConfigs(githubClient, getGenerateConfigsQuery{
		RawFileQuery:  q.RawFileQuery,
		ProgramConfig: p,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get generate configs")
	}

	// Get the schema.
	schema, err := getSchema(githubClient, getSchemaFileQuery{
		RawFileQuery: q.RawFileQuery,
		SchemaFile:   p.SchemaFile,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema file")
	}

	// Copy the program config, but containing the program AST.
	return &Program{
		ProgramConfig:   *p,
		GenerateConfigs: generateConfigs,
		Root:            schema,
	}, nil
}

func parseProgramConfigFile(f []byte, ext string) (*program.ProgramConfig, error) {
	p := program.ProgramConfig{}
	switch ext {
	case ".json":
		if err := json.Unmarshal(f, &p); err != nil {
			return nil, err
		}
	case ".yml", ".yaml":
		if err := yamlio.Unmarshal(f, &p); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(`unsupported file extension, expect ".json", ".yml" or ".yaml", got: %v`, ext)
	}

	if p.SchemaFile == "" {
		p.SchemaFile = "schema.yml"
	}

	return &p, nil
}

// getSchemaFileQuery is the query for getting the schema file.
type getSchemaFileQuery struct {
	githubclient.RawFileQuery
	SchemaFile string
}

// getSchema gets the schema.
func getSchema(githubClient *githubclient.Client, q getSchemaFileQuery) (root *ast.RootNode, err error) {
	f, err := githubClient.RawFile(githubclient.RawFileQuery{
		Owner: q.Owner,
		Repo:  q.Repo,
		Ref:   q.Ref,
		Path:  path.Join(path.Dir(q.Path), q.SchemaFile),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema file from github")
	}

	// Parse the schema file.
	return parseSchemaFile([]byte(f), path.Ext(q.Path))
}

// parseSchemaFile parses the schema file.
func parseSchemaFile(f []byte, ext string) (root *ast.RootNode, err error) {
	switch ext {
	case ".json":
		var s interface{}
		if err := json.Unmarshal(f, &s); err != nil {
			return nil, fmt.Errorf("[Models parsing error]: %v", err)
		}
		b, err := yamlio.Marshal(s)
		if err != nil {
			return nil, fmt.Errorf("[Models parsing error]: %v", err)
		}
		root, err = yaml.GetRootNodeFromYaml(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("[Models parsing error]: %v", err)
		}

	case ".yml", ".yaml":
		root, err = yaml.GetRootNodeFromYaml(bytes.NewReader(f))
		if err != nil {
			return nil, fmt.Errorf("[Models parsing error]: %v", err)
		}

	case ".proto":
		root, err = proto.GetRootNodeFromProto(bytes.NewReader(f))

	default:
		return nil, fmt.Errorf("schema file extension not allowed")
	}

	if err != nil {
		return nil, fmt.Errorf("[Models parsing error]: %v", err)
	}

	return root, nil
}

// getGenerateConfigsQuery is the query for getting the generate configs.
type getGenerateConfigsQuery struct {
	githubclient.RawFileQuery
	ProgramConfig *program.ProgramConfig
}

// getGenerateConfigs gets the generate configs.
func getGenerateConfigs(githubClient *githubclient.Client, q getGenerateConfigsQuery) ([]generators.GenerateConfig, error) {
	var generateConfigs []generators.GenerateConfig
	for _, data := range q.ProgramConfig.Generators {
		tmplFile := path.Join(path.Dir(q.Path), data.Template)
		f, err := githubClient.RawFile(githubclient.RawFileQuery{
			Owner: q.Owner,
			Repo:  q.Repo,
			Ref:   q.Ref,
			Path:  tmplFile,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get template file from github")
		}

		generateConfigs = append(generateConfigs, generators.GenerateConfig{
			Name:    data.Name,
			Output:  data.Output,
			Ignore:  data.Ignore,
			From:    data.From,
			Types:   data.Types,
			Helpers: data.Helpers,

			Template: f,
		}.ApplyDefinitions(q.ProgramConfig.Definitions))
	}

	return generateConfigs, nil
}

// genFiles generates new program files.
func genFiles(root *ast.RootNode, configs []generators.GenerateConfig) ([]generators.OutputFile, error) {
	generated := []generators.OutputFile{}
	for _, g := range configs {
		f, err := generators.Generate(root, g, false)
		if err != nil {
			return nil, err
		}
		generated = append(generated, f...)
	}
	return generated, nil
}
