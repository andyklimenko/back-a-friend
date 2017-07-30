package db

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTournamentsTable = "CREATE TABLE 'Tournaments' (`TourId`	INTEGER NOT NULL UNIQUE, `Deposit`	INTEGER NOT NULL, `Winners`	TEXT, `Players`	TEXT, PRIMARY KEY(TourId));"
	createWinnersTable     = "CREATE TABLE 'Winners' (`Winner-id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `TourId` INTEGER NOT NULL, `Player-id`	TEXT NOT NULL, `Prize` INTEGER NOT NULL);"
	createPlayersTable     = "CREATE TABLE `Players` (`PlayerId` TEXT NOT NULL UNIQUE, `Points`	INTEGER, PRIMARY KEY(PlayerId));"
)

var (
	ErrorNotFound    = errors.New("Not found")
	ErrAlreadyExists = errors.New("Already exists")
)

type Db struct {
	db    *sql.DB
	dbMux sync.Mutex
}

func (d *Db) Create(dbPath string) error {
	var err error
	d.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	var statement *sql.Stmt
	statement, err = tx.Prepare(createTournamentsTable)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(); err != nil {
		return err
	}

	statement, err = tx.Prepare(createWinnersTable)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(); err != nil {
		return err
	}

	statement, err = tx.Prepare(createPlayersTable)
	if err != nil {
		return err
	}
	if _, err = statement.Exec(); err != nil {
		return err
	}
	return tx.Commit()
}

func (d *Db) Stop() error {
	return d.db.Close()
}