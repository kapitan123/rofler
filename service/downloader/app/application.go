package app

import (
	"context"
	"io"

	"github.com/kapitan123/telegrofler/service/downloader/infra"
)

type Application struct {
	videoSavedTopic  successTopic
	videoFilesBucket fileBucket
	downloader       downloader
}

type successTopic interface {
	PublishSuccess(ctx context.Context, savedVideoAddr string, originalUrl string) error
}

type fileBucket interface {
	Save(ctx context.Context, fromReader io.Reader) (string, error)
}

type downloader interface {
	DownloadFromUrl(url string, w io.Writer) error
}

func NewApplicationFromConfig(ctx context.Context, servicename string, projectId string, videoFileBucket string, videoSavedTopicId string) Application {

	videoBucket := infra.NewCloudStoreBucketClient(ctx, projectId, videoFileBucket)
	successTopic := infra.NewPubSubTopicClient(ctx, projectId, servicename, videoSavedTopicId)
	youtubeDl := infra.NewDownloader()

	return NewApplication(successTopic, videoBucket, youtubeDl)
}

func NewApplication(videoSavedTopic successTopic, videoBucket fileBucket, downloader downloader) Application {
	return Application{
		videoSavedTopic:  videoSavedTopic,
		videoFilesBucket: videoBucket,
		downloader:       downloader,
	}
}

func (app *Application) SaveVideoToStorage(ctx context.Context, url string) error {
	pipeReader, pipeWriter := io.Pipe()

	err := app.downloader.DownloadFromUrl(url, pipeWriter)

	if err != nil {
		return err
	}

	id, err := app.videoFilesBucket.Save(ctx, pipeReader)

	if err != nil {
		return err
	}

	err = app.videoSavedTopic.PublishSuccess(ctx, id, url)

	if err != nil {
		return err
	}

	return err
}
