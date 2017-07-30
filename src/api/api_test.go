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

	if err := a.AnnounceTournament(tourId, 1000); err != ErrTournamentAlreadyAnnounced {
		t.Error(errors.New("Created tournament duplicate"))
	}
}

func TestApi_JoinTournamentInsufficientFunds(t *testing.T) {
	const (
		playerId = "P1"
		tourId   = 1
	)

	t.Run("no backers", func(t *testing.T) {
		a, closer, err := setupApi()
		if err != nil {
			t.Fatal(err)
		}
		defer closer()

		if err := a.Fund(playerId, 300); err != nil {
			t.Fatal(err)
		}
		if err := a.AnnounceTournament(tourId, 1000); err != nil {
			t.Fatal(err)
		}
		if err := a.JoinTournament(tourId, playerId, []string{}); err != ErrInsufficientFunds {
			t.Error(err)
		}
	})

	t.Run("3 backers, player with no pts", func(t *testing.T) {
		a, closer, err := setupApi()
		if err != nil {
			t.Fatal(err)
		}
		defer closer()

		if err := a.Fund("P1", 10); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P2", 300); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P3", 300); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P4", 500); err != nil {
			t.Fatal(err)
		}

		if err := a.AnnounceTournament(tourId, 1000); err != nil {
			t.Fatal(err)
		}

		if err := a.JoinTournament(tourId, "P1", []string{"P2", "P3", "P4"}); err != ErrInsufficientFunds {
			t.Fatal("P1 should have no sufficient funds")
		}
	})

	t.Run("3 backers, with not enough pts", func(t *testing.T) {
		a, closer, err := setupApi()
		if err != nil {
			t.Fatal(err)
		}
		defer closer()

		if err := a.Fund("P1", 300); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P2", 30); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P3", 300); err != nil {
			t.Fatal(err)
		}
		if err := a.Fund("P4", 500); err != nil {
			t.Fatal(err)
		}

		if err := a.AnnounceTournament(tourId, 1000); err != nil {
			t.Fatal(err)
		}

		if err := a.JoinTournament(tourId, "P1", []string{"P2", "P3", "P4"}); err != ErrInsufficientFunds {
			t.Fatal("P2 should have no sufficient funds")
		}
	})
}

func TestApi_JoinTournament(t *testing.T) {
	a, closer, err := setupApi()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	if err := a.Fund("P1", 300); err != nil {
		t.Fatal(err)
	}
	if err := a.Fund("P2", 300); err != nil {
		t.Fatal(err)
	}
	if err := a.Fund("P3", 300); err != nil {
		t.Fatal(err)
	}
	if err := a.Fund("P4", 500); err != nil {
		t.Fatal(err)
	}
	if err := a.Fund("P5", 1000); err != nil {
		t.Fatal(err)
	}

	const tourId = 1
	if err := a.AnnounceTournament(tourId, 1000); err != nil {
		t.Fatal(err)
	}

	if err := a.JoinTournament(tourId, "P5", []string{}); err != nil {
		t.Fatal(err)
	}

	if err := a.JoinTournament(tourId, "P1", []string{"P2", "P3", "P4"}); err != nil {
		t.Fatal(err)
	}

	b1, err := a.Balance("P1")
	if err != nil {
		t.Fatal(err)
	}

	b2, err := a.Balance("P2")
	if err != nil {
		t.Fatal(err)
	}

	b3, err := a.Balance("P3")
	if err != nil {
		t.Fatal(err)
	}

	b4, err := a.Balance("P4")
	if err != nil {
		t.Fatal(err)
	}

	b5, err := a.Balance("P5")
	if err != nil {
		t.Fatal(err)
	}

	if b1 != 50 {
		t.Error("wrong ballance", b1)
	}

	if b2 != 50 {
		t.Error("wrong ballance", b2)
	}

	if b3 != 50 {
		t.Error("wrong ballance", b3)
	}

	if b4 != 250 {
		t.Error("wrong ballance", b4)
	}

	if b5 != 0 {
		t.Error("wrong ballance", b5)
	}
}

func TestApi_Take(t *testing.T) {
	a, closer, err := setupApi()
	if err != nil {
		t.Fatal(err)
	}
	defer closer()

	if err := a.Fund("P1", 500); err != nil {
		t.Fatal(err)
	}

	if err := a.Take("P1", 1000); err != ErrInsufficientFunds {
		t.Fatal("Player has too many points")
	}

	if err := a.Take("P1", 300); err != nil {
		t.Fatal(err)
	}

	b, err := a.Balance("P1")
	if err != nil {
		t.Fatal(err)
	}
	if b != 200 {
		t.Error("Wrong player ballance ", b)
	}
}
