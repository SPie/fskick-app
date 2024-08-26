package server

import (
        "net/http"

        "github.com/spie/fskick/internal/seasons"
)

type SeasonsController struct {
        seasonsManager seasons.Manager
}

func NewSeasonsController(seasonsManager seasons.Manager) SeasonsController {
        return SeasonsController{seasonsManager: seasonsManager}
}

func (controller SeasonsController) GetSeasons(res http.ResponseWriter, _ *http.Request) {
        seasons, err := controller.seasonsManager.GetSeasons()
        if err != nil {
                handleInternalServerError(res, err)
                return
        }

        seasonsResponse := make([]seasonResponse, len(seasons))
        for i, season := range seasons {
                seasonsResponse[i] = newSeasonResponseFromSeason(season)
        }

        err = writeJsonResponse(res, map[string][]seasonResponse{"seasons": seasonsResponse})
        if err != nil {
                handleInternalServerError(res, err)
                return
        }
}
