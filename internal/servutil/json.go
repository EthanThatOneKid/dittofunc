package servutil

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes the response as JSON.
func WriteJSON(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
