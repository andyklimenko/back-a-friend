package api

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"api/db"
)

func setupApi() (_ Api, _ func(), rerr error) {
	mydb := &db.Db{}

	tmpDir, err := ioutil.TempDir("", "dbTest")
	if err != nil {
		return nil, func() {}, err
	}
	defer func() {
		if rerr != nil {
			os.Remove(tmpDir)
		}
	}()

	dbDir := path.Join(tmpDir, "db")
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		return nil, func() {}, err
	}

	dbPath := path.Join(dbDir, "testdb.db")
	if err := mydb.Create(dbPath); err != nil {
		return nil, func() {}, err
	}

	defer func() {
		if rerr != nil {
			mydb.Stop()
		}
	}()

	a := &api_impl{db: mydb}
	if err := a.Start(); err != nil {
		return nil, func() {}, err
	}

	return a, func() {
		a.db.Stop()
		os.Remove(tmpDir)
	}, nil
}

func TestApi_Fund(t *testing.T) {
	a, closer, err := setupApi()
	if err != nil {
		t.Fatal(err)
	}

	defer closer()

	players := make(map[string]int)
	players["Joe"] = 100
	players["Bob"] = 500

	for player, pts := range players {
		if err := a.Fund(player, pts); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Fund("Joe", 400); err != nil {
		t.Fatal(err)
	}

	joePts, err := a.Balance("Joe")
	if err != nil {
		t.Error(err)
	}

	bobPts, err := a.Balance("Joe")
	if err != nil {
		t.Error(err)
	}

	if joePts != bobPts {
		t.Error(joePts, bobPts)
	}
}

func TestApi_AnounceTournament(t *testing.T) {
	a, closer, err := setupApi()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	tourId := 42
	if err := a.AnnounceTournament(tourId, 1000); err != nil {
		t.Fatal(err)
	}

	if err := a.AnnounceTournament(tourId, 1000); err != db.ErrAlreadyExists {
		t.Error(errors.New("Created tournament duplicate"))
	}
}
