package app

import (
	"context"
	"io"

	"github.com/kapitan123/telegrofler/service/downloader/infra"
	"github.com/sirupsen/logrus"
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
	DownloadFromUrl(ctx context.Context, url string, w io.Writer) error
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

	err := app.downloader.DownloadFromUrl(ctx, url, pipeWriter)

	logrus.Infof("video piped from %s", url)

	if err != nil {
		logrus.Error(err)
		return err
	}

	id, err := app.videoFilesBucket.Save(ctx, pipeReader)

	logrus.Infof("video saved to bucket %s", id)

	if err != nil {
		logrus.Error(err)
		return err
	}

	err = app.videoSavedTopic.PublishSuccess(ctx, id, url)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return err
}
