package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"errors"
)

const (
	createTournamentsTable = "CREATE TABLE 'Tournaments' ( `TourId`	INTEGER NOT NULL UNIQUE, `Winners` TEXT, `Deposit` INTEGER NOT NULL, `Players` TEXT, PRIMARY KEY(TourId))"
	createWinnersTable = "CREATE TABLE 'Winners' (`Winner-id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `TourId` INTEGER NOT NULL, `Player-id`	TEXT NOT NULL, `Prize` INTEGER NOT NULL)"

	announceTournamentQuery = "insert into Tournaments values (?, '', '', ?)"
	selectTournamentPlayersQuery = "select Players from Tournaments where TourId=?"
)

var (
	ErrorNotFound = errors.New("Not found")
	ErrAlreadyExists = errors.New("Already exists")
)

type Db struct {
	db *sql.DB
}

func (d *Db) Create(dbPath string) error {
	var err error
	d.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	var statement *sql.Stmt
	statement, err = d.db.Prepare(createTournamentsTable)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(); err != nil {
		return err
	}

	statement, err = d.db.Prepare(createWinnersTable)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(); err != nil {
		return err
	}

	return nil
}

func (d *Db) CreateTournament(id int, deposit int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(announceTournamentQuery)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id, deposit)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *Db) JoinTournament(tourId int, playerId string, backers []string) (rerr error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if rerr != nil {
			tx.Rollback()
		}
	}()

	rows, err := tx.Query(selectTournamentPlayersQuery, tourId)
	if err != nil {
		return err
	}

	var players string
	if err := rows.Scan(&players); err != nil {
		return err
	}

	// todo: parse players

	return nil
}
