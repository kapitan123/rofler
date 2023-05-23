package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/kapitan123/telegrofler/service/convertor/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := app.NewApplication(ctx)

	apiRouter := chi.NewRouter()

	setMiddlewares(apiRouter)
	setRoutes(application)

	logrus.Info("Starting HTTP server")

	// AK TODO ADD PORT
	err := http.ListenAndServe(":8080", apiRouter)
	if err != nil {
		logrus.WithError(err).Panic("Unable to start HTTP server")
	}
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.NoCache)
}

// AK TODO i need an httphandlker wrapper for this
// this is where I need codegen
func setRoutes(router *chi.Mux, app app.Application) {
	router.Post("/pubsub/subscriptions/video-published", SaveVideo)
	router.Get("/converted-videos/{videoId}", getVideo)
}
