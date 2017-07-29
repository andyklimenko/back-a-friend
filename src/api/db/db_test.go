package db

import (
	"testing"
	"io/ioutil"
	"os"
	"path"
)

func setupMyDb() (_ *Db, _ func(), rerr error) {
	mydb := &Db{}

	tmpDir, err := ioutil.TempDir("", "dbTest")
	if err != nil {
		return nil, func(){}, err
	}
	defer func() {
		if rerr != nil {
			os.Remove(tmpDir)
		}
	}()

	dbDir := path.Join(tmpDir, "db")
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		return nil, func(){}, err
	}

	dbPath := path.Join(dbDir, "testdb.db")
	if err := mydb.Create(dbPath); err != nil {
		return nil, func(){}, err
	}

	closer := func() {
		mydb.db.Close()
		os.Remove(tmpDir)
	}

	return mydb, closer, nil
}

func TestDb_CreateTournament(t *testing.T) {
	myDb, closer, err := setupMyDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	const tournamentId = 100500
	const deposit = 1000
	if err := myDb.CreateTournament(tournamentId, deposit); err != nil {
		t.Fatal(err)
	}

	rows, err := myDb.db.Query("select * from Tournaments where TourId=?", tournamentId)
	if err != nil {
		t.Fatal(err)
	}

	var tourId, depo int
	var winners, players string

	rows.Next()
	if err := rows.Scan(&tourId, &winners, &players, &depo); err != nil {
		t.Fatal(err)
	}

	if tourId != tournamentId {
		t.Error(tourId)
	}

	if winners != "" {
		t.Error(winners)
	}

	if players != "" {
		t.Error(players)
	}

	if depo != deposit {
		t.Error(depo)
	}
}