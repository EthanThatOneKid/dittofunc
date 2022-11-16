package servutil

import (
	"net/http"
	"strings"
)

// ReadToken reads the token from the request headers.
func ReadToken(r *http.Request, dst *string) {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
	}

	if token != "" {
		*dst = token
	}
}
