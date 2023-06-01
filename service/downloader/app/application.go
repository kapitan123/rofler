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
	DownloadFromUrl(ctx context.Context, url string, out io.Writer) error
}

func NewApplicationFromConfig(ctx context.Context, servicename string, projectId string, videoFileBucket string, videoSavedTopicId string, ytPath string) Application {
	videoBucket := infra.NewCloudStoreBucketClient(ctx, projectId, videoFileBucket)
	successTopic := infra.NewPubSubTopicClient(ctx, projectId, servicename, videoSavedTopicId)
	youtubeDl := infra.NewDownloader(ytPath)

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
	pr, pw := io.Pipe() // AK TODO I can actually change signature to return reader from Download Url
	// this way I

	errs := make(chan error, 1)

	go func() {
		defer pw.Close()
		errs <- app.downloader.DownloadFromUrl(ctx, url, pw)
	}()

	// AK TODO need to peek first byte
	id, err := app.videoFilesBucket.Save(ctx, pr)

	if err := <-errs; err != nil {
		close(errs)
		return err
	}

	logrus.Infof("video saved to bucket %s", id)

	if err != nil {
		return err
	}

	err = app.videoSavedTopic.PublishSuccess(ctx, id, url)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return err
}
