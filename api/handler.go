package handler

import (
	"net/http"

	"github.com/ethanthatonekid/ditto-edge/internal/dittohandler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := dittohandler.Handle(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
