package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api"
)

type announceTournament struct {
	a api.Api
}

func newAnnounceTournament(a api.Api) http.Handler {
	return announceTournament{a}
}

func (h announceTournament) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	tourId, ok := q["tournamentId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deposit, ok := q["deposit"]
	if !ok {
		http.Error(w, "deposit", http.StatusBadRequest)
		return
	}

	if len(tourId) > 1 {
		http.Error(w, "deposit", http.StatusBadRequest)
		return
	}

	tid, err := strconv.Atoi(tourId[0])
	if err != nil {
		http.Error(w, "tourId", http.StatusBadRequest)
		return
	}
	d, err := strconv.Atoi(deposit[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.a.AnnounceTournament(tid, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type joinTournament struct {
	a api.Api
}

func newJoinTournament(a api.Api) http.Handler {
	return joinTournament{a}
}

func (h joinTournament) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	tourId, ok := q["tournamentId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playerId, ok := q["playerId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(tourId) > 1 || len(playerId) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tid, err := strconv.Atoi(tourId[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	backers, ok := q["backerId"]
	if err := h.a.JoinTournament(tid, playerId[0], backers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type resultTournament struct {
	a api.Api
}

func newResultTournament(a api.Api) http.Handler {
	return resultTournament{a}
}

func (h resultTournament) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	winner, err := h.a.ResultTournament()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(winner)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
