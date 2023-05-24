package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/kapitan123/telegrofler/service/downloader/app"
	"github.com/kapitan123/telegrofler/service/downloader/config"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	cfg := config.GetEnvVars()

	logs.Init(cfg.DebguMode)

	application := app.NewApplication(ctx, cfg.ProjectId, cfg.VideoConvertedTopicId, cfg.VidoFilesBucketUrl)

	apiRouter := chi.NewRouter()

	setMiddlewares(apiRouter)
	setRoutes(apiRouter, application)

	logrus.Info("Starting HTTP server")

	err := http.ListenAndServe(cfg.Port, apiRouter)
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
