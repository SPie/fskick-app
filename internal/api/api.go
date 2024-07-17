package api

import (
	"fmt"
	"net/http"
)

type Server interface {
	Run(addr ...string) error
}

type MuxServer struct {
	mux *http.ServeMux
}

func SetUp(
	gamesController GamesController,
	seasonsController SeasonsController,
	playersController PlayersController,
) Server {
	server := MuxServer{mux: http.NewServeMux()}

	server.mux.HandleFunc("/api/seasons", seasonsController.GetSeasons)
	server.mux.HandleFunc("/api/seasons/table", gamesController.GetTable)
	server.mux.HandleFunc("/api/seasons/table/{season}", gamesController.GetTable)

	server.mux.HandleFunc("/api/players", playersController.GetPlayers)
	server.mux.HandleFunc("/api/players/:player/team", playersController.GetFavoriteTeam)
	server.mux.HandleFunc("/api/players/{player}", playersController.GetPlayers)

	server.mux.HandleFunc("/api/games/count", gamesController.GetGamesCount)

	return &server
}

func (server *MuxServer) Run(addr ...string) error {
	fmt.Printf("Starting the server on %s...\n", addr[0])

	return http.ListenAndServe(addr[0], server.mux)
}
