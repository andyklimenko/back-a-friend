package server

import (
	"path"
	"os"

	"api"
	"api/db"
	"net/http"
)

func StartServer(curDir string) (chan struct{}, error) {
	dbDir := path.Join(curDir, "db")
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		return nil, err
	}

	dbPath := path.Join(dbDir, "back-a-friend.db")

	mydb := &db.Db{}
	if err := mydb.Create(dbPath); err != nil {
		return nil, err
	}

	var err error
	a, err := api.CreateApi(mydb)
	if err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})
	go func() {
		// init http-server
		http.Handle("/take", newTakeHandler(a))
		http.Handle("/fund", newFundHandler(a))
		http.Handle("/balance", newBalanceHandler(a))
		http.Handle("/announceTournament", newAnnounceTournament(a))
		http.Handle("/joinTournament", newJoinTournament(a))
		http.Handle("/resultTournament", newResultTournament(a))
		http.ListenAndServe(":8080", nil)
	}()
	return doneCh, nil
}
