package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kapitan123/telegrofler/common/logs"
	"github.com/kapitan123/telegrofler/common/server"
	"github.com/kapitan123/telegrofler/service/downloader/config"
	"github.com/kapitan123/telegrofler/service/downloader/downloader"
	"github.com/kapitan123/telegrofler/service/downloader/port"
)

func main() {
	ctx := context.Background()
	cfg := config.GetEnvVars()

	logs.Init(cfg.DebguMode)

	application := downloader.NewApplicationFromConfig(ctx, cfg.ServiceName, cfg.ProjectId,
		cfg.VidoFilesBucket, cfg.VideoSavedTopic, cfg.YoutubeDlPath, cfg.DownloaderCookies)

	server.RunHTTPServer(cfg.Port, func(router chi.Router) http.Handler {
		return port.HandlerFromMux(port.NewHttpServer(application), router)
	})
}
