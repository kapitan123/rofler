package app

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type Application struct {
	videoSavedTopic  *topic
	videoFilesBucket *fileBucket
	downloader       *downloader
}

type topic interface {
	PublishSuccess(savedVideoAddr string, originalUrl string) error
}

type fileBucket interface {
	Save(w io.Writer) error
	Read(addr string, r io.Reader) error
}

type downloader interface {
	DownloadFromUrl(url string, w io.Writer) error
}

func BoostrapNewApplication(ctx context.Context, projectId string, videoSavedTopicId string, videoFilesBucketUrl string) Application {
	newStorageClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	// I should pass all this shit to adapters and here instaniate wrappers
	// adapter.CloudStorageBucket(ctx, projectId)
	// and then I can pass it as fileBucket interface
	newPubSubClient, err := pubsub.NewClient(ctx, projectId)

	if err != nil {
		panic(err)
	}

	return Application{
		videoSavedTopic:  newPubSubClient.Topic(videoSavedTopicId),
		videoFilesBucket: newStorageClient.Bucket(videoFilesBucketUrl),
		downloader:       youtubeDl.NewDownloader(),
	}
}

func NewApplication(ctx context.Context, videoSavedTopic topic, videoBucket fileBucket) {

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
