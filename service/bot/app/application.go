package app

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/infra"
	"github.com/sirupsen/logrus"
)

// AK TODO looks like this one is a case for CQRS
// top roflers - is a query
// the rest are commands
type Application struct {
	messenger              messenger
	videoFilesBucket       fileBucket
	videoUrlPublishedTopic videoUrlPublishedTopic
	postsStorage           postsStorage
	urlsStorage            urlsStorage
}
type postsStorage interface {
}

type urlsStorage interface {
	GetUrlByAddr(ctx context.Context, url string) (url.Url, bool, error)
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

// AK TODO i think I need to extract telegram commands, and inside handle pass it to the right one
func NewApplication(messenger messenger, postsStorage postsStorage, videoFilesBucket fileBucket, videoUrlPublishedTopic videoUrlPublishedTopic) Application {
	return Application{
		messenger:              messenger,
		videoFilesBucket:       videoFilesBucket,
		videoUrlPublishedTopic: videoUrlPublishedTopic,
		postsStorage:           postsStorage,
	}
}

func (app *Application) PublishVideo(ctx context.Context, originalUrl string, savedAddr string) error {
	// AK TODO implement

	mediaData, found, err := app.urlsStorage.GetUrlByAddr(ctx, originalUrl)

	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("video %s was not found for url %s", savedAddr, originalUrl)
	}
	// get in store - get chat id, get poster, publish video should have a ttl of one month

	return nil
}

func (app *Application) HandleTelegramMessage(ctx context.Context, message tgbotapi.Update) error {

	// AK TODO get command
	// execute command

	// here is the domain logic
	// AK TODO implement
	return nil
}

type MessageHandler func(int)
