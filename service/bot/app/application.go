package app

import (
	"context"
	"io"
	//"github.com/kapitan123/telegrofler/service/downloader/infra"
)

type Application struct {
	messenger              messenger
	videoFilesBucket       fileBucket
	videoUrlPublishedTopic videoUrlPublishedTopic
}

type messenger interface {
}

type videoUrlPublishedTopic interface {
	PublishUrl(ctx context.Context, savedVideoAddr string, originalUrl string) error
}

type fileBucket interface {
	Read(ctx context.Context, addr string, fromReader io.Reader) error
}

// func NewApplicationFromConfig(ctx context.Context, servicename string, projectId string, videoFileBucket string, videoSavedTopicId string) Application {

// 	videoBucket := infra.NewCloudStoreBucketClient(ctx, projectId, videoFileBucket)
// 	successTopic := infra.NewPubSubTopicClient(ctx, projectId, servicename, videoSavedTopicId)
// 	youtubeDl := infra.NewDownloader()

// 	return NewApplication(successTopic, videoBucket, youtubeDl)
// }

func NewApplication(messenger messenger, videoFilesBucket fileBucket, videoUrlPublishedTopic videoUrlPublishedTopic) Application {
	return Application{
		messenger:              messenger,
		videoFilesBucket:       videoFilesBucket,
		videoUrlPublishedTopic: videoUrlPublishedTopic,
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

func (app *Application) GetVideo() {
	// AK TODO implement
}
