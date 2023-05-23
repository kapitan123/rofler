package app

import (
	"context"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type Application struct {
	videoConvertedTopic *pubsub.Topic
	videFilesBucket     *storage.BucketHandle
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
	}
}

func (app *Application) SaveVideoToStorage() {

}

// These wrappers can be put in a separate repo, but it seem like an overkill
func (app *Application) GetVideo() {
	// if it is not present tell the link is expired
}

func (app *Application) save() {

}

func (app *Application) download() {

}
