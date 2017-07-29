package db

import "strings"

const (
	announceTournamentQuery      = "insert into Tournaments values (?, ?, '', '')"
	selectTournamentPlayersQuery = "select Players from Tournaments where TourId=?"
	updateTournamentPlayersQuery = "update Tournaments set Players=? where TourId = ?"
)

func (d *Db) CreateTournament(id int, deposit int) error {
	d.dbMux.Lock()
	defer d.dbMux.Unlock()

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

func (d *Db) JoinTournament(tourId int, playerId string) (rerr error) {
	d.dbMux.Lock()
	defer d.dbMux.Unlock()

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

	if !rows.Next() {
		return ErrorNotFound
	}

	var players string
	if err := rows.Scan(&players); err != nil {
		return err
	}

	pArr := strings.Split(players, ",")
	for _, p := range pArr {
		if p == playerId {
			return ErrAlreadyExists
		}
	}

	semicolonRequired := len(pArr) >= 1 && pArr[0] != ""
	if semicolonRequired {
		players += "," + playerId
	} else {
		players += playerId
	}

	stmt, err := tx.Prepare(updateTournamentPlayersQuery)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(players, tourId); err != nil {
		return err
	}

	return tx.Commit()
}
