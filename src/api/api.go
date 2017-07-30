package api

import (
	"api/db"
)

type Winner struct {
	PlayerId string
	Prize    int
}

type Api interface {
	Start() error
	Take(playerToTakeFrom string, points int) error
	Fund(playerId string, points int) error
	AnnounceTournament(tourId int, deposit int) error
	JoinTournament(tourId, playerId string) error
	ResultTournament() ([]Winner, error)
	Balance(playerId string) (int, error)
}

type api_impl struct {
	db                 *db.Db
	activeTournamentId int
}

func (a *api_impl) Start() error {
	return nil
}

func (a *api_impl) Take(playerToTakeFrom string, points int) error {
	return nil
}

func (a *api_impl) Fund(playerId string, points int) error {
	pts, err := a.db.PlayerPoints(playerId)
	if err == nil {
		//update
		if err := a.db.UpdatePlayer(playerId, pts+points); err != nil {
			return err
		}
		return nil
	} else if err == db.ErrorNotFound {
		//add new
		if err := a.db.CreatePlayer(playerId, points); err != nil {
			return err
		}
		return nil
	}
	return err
}

func (a *api_impl) AnnounceTournament(tourId int, deposit int) error {
	exists, err := a.db.TournamentExists(tourId)
	if err != nil {
		return err
	}
	if exists {
		return db.ErrAlreadyExists
	}

	return a.db.CreateTournament(tourId, deposit)
}

func (a *api_impl) JoinTournament(tourId, playerId string) error {
	//todo: implement
	return nil
}

func (a *api_impl) ResultTournament() ([]Winner, error) {
	//todo: implement
	return []Winner{}, nil
}

func (a *api_impl) Balance(playerId string) (int, error) {
	return a.db.PlayerPoints(playerId)
}
