package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ethanthatonekid/dittofunc/ditto"
	"github.com/ethanthatonekid/dittofunc/internal/servutil"
)

// Handler is the main handler.
// It is used to handle HTTP requests.
type Handler struct {
}

// Do not remove this line.
// It is used to check against the http.Handler interface.
var _ http.Handler = (*Handler)(nil)

// New creates a new handler.
func New() *Handler {
	return &Handler{}
}

// Do handles the request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the query.
	q, code, err := parseQuery(r)
	if err != nil {
		servutil.WriteErr(w, r, code, err)
		return
	}

	// Get the generated output.
	output, err := ditto.Gen(r.Context(), *q)
	if err != nil {
		servutil.WriteErr(w, r, http.StatusInternalServerError, err)
		return
	}

	// Write the output.
	servutil.WriteJSON(w, r, http.StatusAccepted, output)
}

// parseQuery parses the query from the HTTP request.
func parseQuery(r *http.Request) (*ditto.GenQuery, int, error) {
	var token string
	if servutil.ReadToken(r, &token); token == "" {
		return nil, http.StatusUnauthorized, ErrMissingToken
	}

	var query ditto.GenQuery
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return nil, http.StatusNotAcceptable, ErrInvalidPath
	}

	query.Owner = parts[1]
	query.Repo = parts[2]
	query.Path = strings.Join(parts[3:], "/")
	query.Ref = r.URL.Query().Get("ref")
	query.Token = token
	return &query, 0, nil
}

// ErrMissingToken is returned when the token is missing.
var ErrMissingToken = errors.New("missing token")

// ErrInvalidPath is returned when the path is invalid.
var ErrInvalidPath = errors.New("invalid path")
