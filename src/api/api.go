package api

type Winner struct {
	PlayerId string
	Prize int
}

type Api interface {
	Take(playerId string, pts int) error
	Fund(playerId string, pts int) error
	AnnounceTournament(tourId string, deposit int) error
	JoinTournament(tourId, playerId string) error
	ResultTournament() ([]Winner, error)
}