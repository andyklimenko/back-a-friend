package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTournamentsTable = "CREATE TABLE 'Tournaments' (`TourId`	INTEGER NOT NULL UNIQUE, `Deposit`	INTEGER NOT NULL, `Players`	TEXT, PRIMARY KEY(TourId));"
	createPlayersTable     = "CREATE TABLE `Players` (`PlayerId` TEXT NOT NULL UNIQUE, `Points`	INTEGER, PRIMARY KEY(PlayerId));"

	deleteTournamentsQuery = "DELETE FROM Tournaments;"
	deletePlayersQuery     = "DELETE FROM Players;"
)

var (
	ErrorNotFound    = errors.New("Not found")
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

func (d *Db) Reset() (rerr error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rerr != nil {
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec(deleteTournamentsQuery); err != nil {
		return err
	}

	if _, err := tx.Exec(deletePlayersQuery); err != nil {
		return err
	}
	return tx.Commit()
}
