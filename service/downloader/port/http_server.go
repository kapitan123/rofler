// it is a wrapper around application to format requests
package port

import (
	"net/http"

	"github.com/kapitan123/telegrofler/service/downloader/app"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app: app}
}

func (h HttpServer) GetTrainings(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var appTrainings []query.Training

	if user.Role == "trainer" {
		appTrainings, err = h.app.Queries.AllTrainings.Handle(r.Context(), query.AllTrainings{})
	} else {
		appTrainings, err = h.app.Queries.TrainingsForUser.Handle(r.Context(), query.TrainingsForUser{User: user})
	}

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	trainings := appTrainingsToResponse(appTrainings)
	trainingsResp := Trainings{trainings}

	render.Respond(w, r, trainingsResp)
}
