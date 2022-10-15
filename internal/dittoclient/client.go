package dittoclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juanvillacortac/ditto/pkg/program"
)

// TODO: Parse the following data from the request:
// - [ProgramConfig](https://github.com/juanvillacortac/ditto/blob/7118ccafe58f218892627b6dd1eb4601781b591b/pkg/program/program.go#L18)
// - [GenerateConfig](https://github.com/juanvillacortac/ditto/blob/7118ccafe58f218892627b6dd1eb4601781b591b/pkg/program/program.go#L27)
// - Or parse the request body into a ProgramConfig and GenerateConfig; [Parse](https://github.com/juanvillacortac/ditto/blob/7118ccafe58f218892627b6dd1eb4601781b591b/pkg/program/program.go#L61)

func ParseConfigFromRequest(r *http.Request) (*program.ProgramConfig, *program.GenerateConfig, error) {
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
