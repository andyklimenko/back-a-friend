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

	handlers := make(map[string]http.Handler)
	handlers["/take"] = newTakeHandler(a)
	handlers["/fund"] = newFundHandler(a)
	handlers["/balance"] = newBalanceHandler(a)

	doneCh := make(chan struct{})
	go func() {
		// init http-server
		http.Handle("/take", handlers["/take"])
		http.Handle("/fund", handlers["/fund"])
		http.Handle("/balance", handlers["/balance"])
		http.ListenAndServe(":8080", nil)
		close(doneCh)
	}()
	return doneCh, nil
}
