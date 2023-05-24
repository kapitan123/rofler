package app

import (
	"context"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/kapitan123/telegrofler/service/downloader/adapter/youtubeDl"
)

type Application struct {
	videoConvertedTopic *pubsub.Topic // AK TODO I suppose to abstract it away I guess
	videFilesBucket     *storage.BucketHandle
	downloader          *youtubeDl.Downloader
}

func NewApplication(ctx context.Context, projectId string, videoConvertedTopicId string, videFilesBucketUrl string) Application {
	newStorageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	newPubSubClient, err := pubsub.NewClient(ctx, projectId)

	if err != nil {
		panic(err)
	}

	return Application{
		videoConvertedTopic: newPubSubClient.Topic(videoConvertedTopicId),
		videFilesBucket:     newStorageClient.Bucket(videFilesBucketUrl),
		downloader:          youtubeDl.NewDownloader(),
	}
}

func (app *Application) SaveVideoToStorage() {
	// wc := d.bucket.Object(fileName).NewWriter(d.ctx)
}

// These wrappers can be put in a separate repo, but it seem like an overkill
func (app *Application) GetVideo() {
	// if it is not present tell the link is expired
}

func (app *Application) save() {

}

func (app *Application) download() {

}
