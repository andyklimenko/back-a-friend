package api

import (
	"errors"
	"sync"

	"api/db"
)

var (
	ErrInsufficientFunds  = errors.New("Not enough funds")
	ErrInvalidQueryResult = errors.New("Invalid query result")
)

type Winner struct {
	PlayerId string
	Prize    int
}

type funded struct {
	backers   []string
	fundedPts int
}

type Api interface {
	Start() error
	Take(playerToTakeFrom string, points int) error
	Fund(playerId string, points int) error
	AnnounceTournament(tourId int, deposit int) error
	JoinTournament(tourId int, playerId string, backers []string) error
	ResultTournament() ([]Winner, error)
	Balance(playerId string) (int, error)
}

type api_impl struct {
	db                 *db.Db
	activeTournamentId int
	funded             map[string]*funded
	dbMux              sync.Mutex
}

func (a *api_impl) Start() error {
	a.funded = make(map[string]*funded)
	return nil
}

func (a *api_impl) Take(playerToTakeFrom string, points int) error {
	return nil
}

func (a *api_impl) Fund(playerId string, points int) error {
	a.dbMux.Lock()
	defer a.dbMux.Unlock()

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
	a.dbMux.Lock()
	defer a.dbMux.Unlock()

	info, err := a.db.TournamentInfo(tourId)
	if err != nil && err != db.ErrorNotFound {
		return err
	}
	if info != nil {
		return db.ErrAlreadyExists
	}

	return a.db.CreateTournament(tourId, deposit)
}

func (a *api_impl) JoinTournament(tourId int, playerId string, backers []string) error {
	a.dbMux.Lock()
	defer a.dbMux.Unlock()

	info, err := a.db.TournamentInfo(tourId)
	if err != nil {
		return err
	}

	balance, err := a.balance(playerId)
	if err != nil {
		return err
	}

	if len(backers) == 0 && info.Deposit > balance {
		return ErrInsufficientFunds
	}

	if balance >= info.Deposit && len(backers) == 0 {
		return a.db.JoinTournament(tourId, playerId)
	}

	requiredPts := info.Deposit / (len(backers) + 1) // backers + playerId
	if balance < requiredPts {
		return ErrInsufficientFunds
	}

	backersMap, err := a.db.MultiplePlayerPoints(backers)
	if err != nil {
		return err
	}

	if len(backersMap) != len(backers) {
		return ErrInvalidQueryResult
	}

	for _, pts := range backersMap {
		if pts <= requiredPts {
			return ErrInsufficientFunds
		}
	}

	f := &funded{[]string{}, requiredPts}
	for b, pts := range backersMap {
		if err := a.db.UpdatePlayer(b, pts-requiredPts); err != nil {
			return err
		}
		f.backers = append(f.backers, b)
	}
	a.funded[playerId] = f

	return a.db.JoinTournament(tourId, playerId)
}

func (a *api_impl) ResultTournament() ([]Winner, error) {
	//todo: implement
	return []Winner{}, nil
}

func (a *api_impl) Balance(playerId string) (int, error) {
	a.dbMux.Lock()
	defer a.dbMux.Unlock()

	return a.balance(playerId)
}

func (a *api_impl) balance(playerId string) (int, error) {
	return a.db.PlayerPoints(playerId)
}
