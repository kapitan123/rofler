package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kapitan123/telegrofler/service/downloader/app"
	"github.com/kapitan123/telegrofler/service/downloader/downloader"
)

func main() {
	ctx := context.Background()
	cfg := app.GetEnvVars()

	// AK TODO add errors
	videoBucket := downloader.NewCloudStoreBucketClient(ctx, cfg.ProjectId, cfg.VideoFilesBucket)
	successTopic := downloader.NewPubSubTopicClient(ctx, cfg.ProjectId, cfg.ServiceName, cfg.VideoSavedTopicId)
	youtubeDl := downloader.NewDownloader(cfg.DownloaderCookies)

	handler := app.NewServer(videoBucket,
		youtubeDl,
		successTopic,
	)

	http := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	go func() {
		log.
	}
}
