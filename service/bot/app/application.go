package app

import (
	"context"
	"io"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/infra"
	"github.com/sirupsen/logrus"
)

type Application struct {
	messenger              messenger
	videoFilesBucket       fileBucket
	videoUrlPublishedTopic videoUrlPublishedTopic
	postsStorage           postsStorage
}
type postsStorage interface {
}

type messenger interface {
}

type videoUrlPublishedTopic interface {
	PublishUrl(ctx context.Context, url string) error
}

type fileBucket interface {
	Read(ctx context.Context, addr string, fromReader io.Reader) error
}

func NewApplicationFromConfig(ctx context.Context, servicename string, projectId string, telegramBotToken string, videoUrlPostedTopicId string, videoFilesBucketUrl string) Application {
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		logrus.Panic("failed to create firestore client")
	}

	postsRepo := infra.NewFirestorePostsRepository(client)

	botapi := infra.NewMessenger(telegramBotToken)

	fileBucket := infra.NewCloudStoreBucketClient(ctx, projectId, videoFilesBucketUrl)
	urlPostedTopic := infra.NewPubSubTopicClient(ctx, projectId, servicename, videoUrlPostedTopicId)

	return NewApplication(botapi, postsRepo, fileBucket, urlPostedTopic)
}

func NewApplication(messenger messenger, postsStorage postsStorage, videoFilesBucket fileBucket, videoUrlPublishedTopic videoUrlPublishedTopic) Application {
	return Application{
		messenger:              messenger,
		videoFilesBucket:       videoFilesBucket,
		videoUrlPublishedTopic: videoUrlPublishedTopic,
		postsStorage:           postsStorage,
	}
}

// func (app *Application) RequestVideoDownload(ctx context.Context, url string) error {
// 	return nil
// }

//	func (app *Application) GetTopRoflers() {
//		// AK TODO implement
//	}
func (app *Application) PublishVideo(ctx context.Context, originalUrl string, savedAddr string) error {
	// AK TODO implement
	return nil
}

func (app *Application) GetCommand(ctx context.Context, message tgbotapi.Update) MessageHandler {

	// here is the domain logic
	// AK TODO implement
	return nil
}

type MessageHandler func(int)
