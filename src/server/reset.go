package server

import (
	"net/http"

	"api"
)

type resetHandler struct {
	a api.Api
}

func newResetHandler(a api.Api) http.Handler {
	return resetHandler{a}
}

func (h resetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.a.Reset()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
