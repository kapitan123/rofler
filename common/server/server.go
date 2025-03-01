package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/sirupsen/logrus"
)

func RunHTTPServer(port int, createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	logrus.Infof("Starting HTTP server on :" + strconv.Itoa(port))

	err := http.ListenAndServe(":"+strconv.Itoa(port), createHandler(apiRouter))
	if err != nil {
		logrus.WithError(err).Panic("Unable to start HTTP server")
	}
}
func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.NoCache)
}
