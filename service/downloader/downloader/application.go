package downloader

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
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

func NewApplicationFromConfig(ctx context.Context, servicename string, projectId string, videoFileBucket string, videoSavedTopicId string, ytPath string, base64cookies string) Application {
	videoBucket := NewCloudStoreBucketClient(ctx, projectId, videoFileBucket)
	successTopic := NewPubSubTopicClient(ctx, projectId, servicename, videoSavedTopicId)
	youtubeDl := NewDownloader(ytPath, base64cookies)

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
	pr, pw := io.Pipe()

	errs := make(chan error, 1)

	go func() {
		defer pw.Close()
		errs <- app.downloader.DownloadFromUrl(ctx, url, pw)
	}()

	bufReader := bufio.NewReader(pr)

	_, err := bufReader.Peek(1)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't download file from url %s", url))
	}

	id, err := app.videoFilesBucket.Save(ctx, bufReader)

	if err := <-errs; err != nil {
		close(errs)
		return err
	}

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't save stream content of %s to bucket %s", url, id))
	}

	logrus.Infof("video saved to bucket %s", id)

	err = app.videoSavedTopic.PublishSuccess(ctx, id, url)

	if err != nil {
		logrus.Error(err)
		return errors.Wrap(err, fmt.Sprintf("can't publish success for saved file %s", id))
	}

	return err
}
