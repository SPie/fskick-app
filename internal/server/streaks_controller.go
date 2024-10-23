package server

import (
	"net/http"

	"github.com/spie/fskick/internal/streaks"
	"github.com/spie/fskick/internal/views"
)

type StreaksViews struct {
	StreaksPage views.StreaksPage
	CurrentStreaks views.CurrentStreaks
}

func NewStreaksViews() StreaksViews {
	return StreaksViews{}
}

type StreaksController struct {
	streaksManager streaks.Manager
	views StreaksViews
}

func NewStreaksController(streaksManager streaks.Manager, views StreaksViews) StreaksController {
	return StreaksController{
		streaksManager: streaksManager,
		views: views,
	}
}

func (controller StreaksController) StreaksPage(res http.ResponseWriter, req *http.Request) {
	longestWiningStreak, longestLosingStreak, err := controller.streaksManager.
		GetLongestWinningAndLosingStreaks()
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	currentStreaks, err := controller.streaksManager.GetCurrentStreaks(true)
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.StreaksPage.Render(
		longestWiningStreak,
		longestLosingStreak,
		currentStreaks,
		req.Context(),
		res,
	); err != nil {
		handleInternalServerError(res, err)
		return
	}
}

func (controller StreaksController) CurrentStreaks(res http.ResponseWriter, req *http.Request) {
	currentStreaks, err := controller.streaksManager.GetCurrentStreaks(req.URL.Query().Get("win") == "on")
	if err != nil {
		handleInternalServerError(res, err)
		return
	}

	if err = controller.views.CurrentStreaks.Render(currentStreaks, req.Context(), res); err != nil {
		handleInternalServerError(res, err)
		return
	}
}
