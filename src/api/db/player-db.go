package db

import "strings"

const (
	playerPtsGetQuery = "Select Points from Players Where PlayerId = ?"
	playerCreateQuery = "Insert into Players values (?, ?)"
	playerUpdateQuery = "Update Players SET Points = ? WHERE PlayerId = ?"
)

func getMultiplePlayersQuery(ids []string) string {
	return "select * from Players where PlayerId in (?" + strings.Repeat(",?", len(ids)-1) + ") order by PlayerId"
}

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

func (d *Db) MultiplePlayerPoints(playerIds []string) (_ map[string]int, rerr error) {
	qry := getMultiplePlayersQuery(playerIds)
	args := []interface{}{}
	for _, id := range playerIds {
		args = append(args, id)
	}

	rows, err := d.db.Query(qry, args...)
	if err != nil {
		return nil, err
	}

	var pid string
	var pts int
	res := make(map[string]int)
	for rows.Next() {
		if err := rows.Scan(&pid, &pts); err != nil {
			return nil, err
		}
		res[pid] = pts
	}

	return res, nil
}

func (d *Db) UpdatePlayer(pid string, pts int) (rerr error) {
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
