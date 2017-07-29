package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTournamentsTable = "CREATE TABLE 'Tournaments' ( `TourId` INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, `Winners` TEXT NOT NULL )"
	createWinnersTable = "CREATE TABLE `Winners` ( `Winner-id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `Player-id` INTEGER NOT NULL, `Prize` INTEGER NOT NULL )"
)

type Db struct {
	db *sql.DB
}

func (d Db) Create(dbPath string) error {
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