package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api"
)

type takeHandler struct {
	a api.Api
}

func newTakeHandler(a api.Api) http.Handler {
	return takeHandler{a}
}

func (h takeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	playerId, ok := q["playerId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pts, ok := q["points"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(pts) > 1 || len(playerId) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := strconv.Atoi(pts[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.a.Take(playerId[0], p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type fundHandler struct {
	a api.Api
}

func newFundHandler(a api.Api) http.Handler {
	return fundHandler{a}
}

func (h fundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	playerId, ok := q["playerId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pts, ok := q["points"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(pts) > 1 || len(playerId) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := strconv.Atoi(pts[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.a.Fund(playerId[0], p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type balanceHandler struct {
	a api.Api
}

func newBalanceHandler(a api.Api) http.Handler {
	return balanceHandler{a}
}

type Balance struct {
	PlayerId string
	Balance  int
}

func (h balanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	playerId, ok := q["playerId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(playerId) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := h.a.Balance(playerId[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := Balance{playerId[0], balance}
	js, err := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
