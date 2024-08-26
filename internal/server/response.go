package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/seasons"
)

type seasonResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newSeasonResponseFromSeason(season seasons.Season) seasonResponse {
	return seasonResponse{
		UUID: season.UUID,
		Name: season.Name,
		Active: season.Active,
		CreatedAt: season.CreatedAt,
		UpdatedAt: season.UpdatedAt,
	}
}

type seasonWithGamesCountResponse struct {
	seasonResponse
	GamesCount int `json:"gamesCount"`
}

func newSeasonsWithGamesCountResponse(season seasons.Season, gamesCount int) seasonWithGamesCountResponse {
	return seasonWithGamesCountResponse{
		seasonResponse: newSeasonResponseFromSeason(season),
		GamesCount: gamesCount,
	}
}

type playerStatsResponses []playerStatsResponse

func newPlayerStatsResponsesFromPlayerStats(playerStats []games.PlayerStats) playerStatsResponses {
	playerStatsRes := make(playerStatsResponses, len(playerStats))
	for i, playerStat := range playerStats {
		playerStatsRes[i] = newPlayerStatsResponseFromPlayerStats(playerStat)
	}

	return playerStatsRes
}

type playerStatsResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Wins int `json:"wins"`
	Games int `json:"games"`
	GamesRatio float64 `json:"gamesRatio"`
	PointsRatio float64 `json:"pointsRatio"`
	Points int `json:"points"`
	WinRatio float64 `json:"winRatio"`
	Position int `json:"position"`
}

func newPlayerStatsResponseFromPlayerStats(playerStats games.PlayerStats) playerStatsResponse {
	return playerStatsResponse{
		UUID: playerStats.UUID,
		Name: playerStats.Name,
		CreatedAt: playerStats.CreatedAt,
		UpdatedAt: playerStats.UpdatedAt,
		Wins: playerStats.Wins,
		Games: playerStats.Games,
		GamesRatio: playerStats.GamesRatio,
		PointsRatio: playerStats.PointsRatio,
		Points: playerStats.Points,
		WinRatio: playerStats.WinRatio,
		Position: playerStats.Position,
	}
}

type tableResponse struct {
	Season seasonWithGamesCountResponse `json:"season"`
	PlayerStats []playerStatsResponse `json:"playerStats"`
}

func newTableResponse(season seasons.Season, gamesCount int, playerStats []games.PlayerStats) tableResponse {
	return tableResponse{
		Season: newSeasonsWithGamesCountResponse(season, gamesCount),
		PlayerStats: newPlayerStatsResponsesFromPlayerStats(playerStats),
	}
}

func writeJsonResponse(res http.ResponseWriter, response interface{}) error {
	jsonRes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("write json response: %w", err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonRes)

	return nil
}

func handleInternalServerError(res http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(res, "Something went wrong.", http.StatusInternalServerError)
}
