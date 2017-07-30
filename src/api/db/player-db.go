package db

const (
	playerPtsGetQuery = "Select Points from Players Where PlayerId = ?"
	playerCreateQuery = "Insert into Players values (?, ?)"
	playerUpdateQuery = "Update Players SET Points = ? WHERE PlayerId = ?"
)

func (d *Db) PlayerPoints(playerId string) (_ int, rerr error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if rerr != nil {
			tx.Rollback()
		}
	}()

	rows, err := tx.Query(playerPtsGetQuery, playerId)
	if err != nil {
		return 0, err
	}

	var pts int
	if !rows.Next() {
		return 0, ErrorNotFound
	}

	if err := rows.Scan(&pts); err != nil {
		return 0, err
	}

	return pts, tx.Commit()
}

func (d *Db) UpdatePlayer(pid string, pts int) (rerr error) {
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

	stmt, err := tx.Prepare(playerUpdateQuery)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(pts, pid); err != nil {
		return err
	}
	return tx.Commit()
}

func (d *Db) CreatePlayer(pid string, pts int) (rerr error) {
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

	stmt, err := tx.Prepare(playerCreateQuery)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(pid, pts); err != nil {
		return err
	}
	return tx.Commit()
}
