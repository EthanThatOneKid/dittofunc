package servutil

import (
	"log"
	"net/http"
)

// WriteErr writes the error as JSON.
func WriteErr(w http.ResponseWriter, _ *http.Request, code int, err error) {
	var errBody struct {
		Error string `json:"error"`
	}

	if err != nil {
		errBody.Error = err.Error()
		log.Println("request", code, "error:", err)
	} else {
		errBody.Error = http.StatusText(code)
		log.Println("request", code, "error")
	}

	WriteJSON(w, nil, code, errBody)
}
