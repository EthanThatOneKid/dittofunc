package handler

import (
	"fmt"
	"net/http"

	"github.com/ethanthatonekid/ditto-edge/internal/dittoclient"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a ProgramConfig and GenerateConfig
	programConfig, generateConfig, err := dittoclient.ParseConfigFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}
