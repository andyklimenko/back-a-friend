package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
		mydb.Stop()
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
	var p string

	i := 0
	for rows.Next() {
		if i > 0 {
			t.Fatal("Too many records")
		}
		if err := rows.Scan(&tourId, &depo, &p); err != nil {
			t.Fatal(err)
		}

		if tourId != tournamentId {
			t.Error(tourId)
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
	playerId1 := "P1"
	playerId2 := "P2"

	if err != myDb.CreatePlayer(playerId1, deposit) {
		t.Error(err)
	}
	if err != myDb.CreatePlayer(playerId2, deposit) {
		t.Error(err)
	}

	if err := myDb.CreateTournament(tournamentId, deposit); err != nil {
		t.Fatal(err)
	}
	if err := myDb.JoinTournament(tournamentId+1, playerId1); err != ErrorNotFound {
		t.Error(errors.New("Unexpected tournament found"))
	}

	if err := myDb.JoinTournament(tournamentId, playerId1); err != nil {
		t.Fatal(err)
	}
	if err := myDb.JoinTournament(tournamentId, playerId2); err != nil {
		t.Fatal(err)
	}

	if unexpected, err := myDb.TournamentInfo(tournamentId + 1); err != ErrorNotFound {
		t.Error(errors.New(fmt.Sprintf("Unexpected tournament found %v", unexpected)))
	}

	tournament, err := myDb.TournamentInfo(tournamentId)
	if err != nil {
		t.Fatal(err)
	}

	if tournament.Id != tournamentId {
		t.Error(tournament.Id)
	}

	if len(tournament.Players) != 2 {
		t.Error(tournament.Players)
	}

	if tournament.Players[0] != playerId1 && tournament.Players[1] != playerId2 {
		t.Error(tournament.Players)
	}

	if tournament.Deposit != deposit {
		t.Error(tournament.Deposit)
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

func TestDb_MultiplePlayerPoints(t *testing.T) {
	myDb, closer, err := setupMyDb()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	const deposit = 1000
	playerId1 := "P1"
	playerId2 := "P2"

	if err != myDb.CreatePlayer(playerId1, deposit) {
		t.Fatal(err)
	}
	if err != myDb.CreatePlayer(playerId2, deposit) {
		t.Fatal(err)
	}

	players, err := myDb.MultiplePlayerPoints([]string{playerId2, playerId1})
	if err != nil {
		t.Fatal(err)
	}

	if len(players) != 2 {
		t.Fatal("wrong players number ", len(players))
	}

	p1Depo, ok := players[playerId1]
	if !ok {
		t.Error(fmt.Sprintf("%s can be found", playerId1))
	}
	if p1Depo != deposit {
		t.Error("wrond deposit ", p1Depo)
	}

	p2Depo, ok := players[playerId2]
	if !ok {
		t.Error(fmt.Sprintf("%s can be found", playerId2))
	}
	if p2Depo != deposit {
		t.Error("wrond deposit ", p2Depo)
	}
}
