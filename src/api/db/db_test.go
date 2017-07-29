package db

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func setupMyDb() (_ *Db, _ func(), rerr error) {
	mydb := &Db{}

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

	closer := func() {
		mydb.db.Close()
		os.RemoveAll(tmpDir)
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
	var w, p string

	i := 0
	for rows.Next() {
		if i > 0 {
			t.Fatal("Too many records")
		}
		if err := rows.Scan(&tourId, &depo, &w, &p); err != nil {
			t.Fatal(err)
		}

		if tourId != tournamentId {
			t.Error(tourId)
		}

		if w != "" {
			t.Error(w)
		}

		if p != "" {
			t.Error(p)
		}

		if depo != deposit {
			t.Error(depo)
		}
		i++
	}
}

func TestDb_JoinTournament(t *testing.T) {
	myDb, closer, err := setupMyDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	const tournamentId = 42
	const deposit = 1000
	if err := myDb.CreateTournament(tournamentId, deposit); err != nil {
		t.Fatal(err)
	}

	playerId1 := "P1"
	playerId2 := "P2"
	if err := myDb.JoinTournament(tournamentId+1, playerId1); err != ErrorNotFound {
		t.Error(err)
	}

	if err := myDb.JoinTournament(tournamentId, playerId1); err != nil {
		t.Fatal(err)
	}
	if err := myDb.JoinTournament(tournamentId, playerId2); err != nil {
		t.Fatal(err)
	}

	rows, err := myDb.db.Query("select * from Tournaments where TourId=?", tournamentId)
	if err != nil {
		t.Fatal(err)
	}

	var tourId, depo int
	var w, p string

	players := []string{}
	players = append(players, playerId1)
	players = append(players, playerId2)

	for rows.Next() {
		if err := rows.Scan(&tourId, &depo, &w, &p); err != nil {
			t.Fatal(err)
		}

		if tourId != tournamentId {
			t.Error(tourId)
		}

		if w != "" {
			t.Error(w)
		}

		plrs := strings.Split(p, ",")
		for i, player := range plrs {
			if player != players[i] {
				t.Error(plrs)
			}
		}

		if depo != deposit {
			t.Error(depo)
		}
	}

	if err := myDb.JoinTournament(tournamentId, playerId1); err != ErrAlreadyExists {
		t.Error(err)
	}
}

func TestDb_PlayerPointsNegative(t *testing.T) {
	myDb, closer, err := setupMyDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	if _, err := myDb.PlayerPoints("no player created yet"); err != ErrorNotFound {
		t.Error(err)
	}
}

func TestDb_FundPlayer(t *testing.T) {
	myDb, closer, err := setupMyDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	players := make(map[string]int)
	players["Joe"] = 100
	players["Bob"] = 500

	for player, pts := range players {
		if err := myDb.FundPlayer(player, pts); err != nil {
			t.Fatal(err)
		}
	}

	if err := myDb.FundPlayer("Joe", 400); err != nil {
		t.Fatal(err)
	}

	joePts, err := myDb.PlayerPoints("Joe")
	if err != nil {
		t.Error(err)
	}

	bobPts, err := myDb.PlayerPoints("Joe")
	if err != nil {
		t.Error(err)
	}

	if joePts != bobPts {
		t.Error(joePts, bobPts)
	}
}
