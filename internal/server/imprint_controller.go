package server

import (
    "net/http"

    "github.com/spie/fskick/internal/views"
)

type ImprintController struct {
    imprintText string
    imprintView views.ImprintView
}

func NewImprintController(imprintText string, imprintView views.ImprintView) ImprintController {
    return ImprintController{
        imprintText: imprintText,
        imprintView: imprintView,
    }
}

func (controller ImprintController) Imprint(res http.ResponseWriter, req *http.Request) {
    err := controller.imprintView.Render(controller.imprintText, req.Context(), res)
    if err != nil {
        handleInternalServerError(res, err)
        return
    }
}
