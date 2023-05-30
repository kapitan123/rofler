package app

import (
	"context"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kapitan123/telegrofler/service/bot/app/pubsubcommand"
	"github.com/kapitan123/telegrofler/service/bot/app/tgcommand"
	"github.com/kapitan123/telegrofler/service/bot/domain/message"
	"github.com/kapitan123/telegrofler/service/bot/infra"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Application struct {
	PublishDownloadedVideo *pubsubcommand.PublishDownloadedVideo
	TelegramCommands       []tgcommandhandler
}

type tgcommandhandler interface {
	Handle(ctx context.Context, message message.Message) error
	ShouldRun(m message.Message) bool
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

	commands := []tgcommandhandler{
		tgcommand.NewRecordMediaPost(postsRepo),
		tgcommand.NewRecordReaction(postsRepo),
		tgcommand.NewRecordUrl(botapi, postsRepo, urlPostedTopic),
		tgcommand.NewTopRofler(botapi, postsRepo),
	}

	publishVideo := pubsubcommand.NewPublishDownloadedVideo(botapi, postsRepo, fileBucket)

	return Application{
		PublishDownloadedVideo: publishVideo,
		TelegramCommands:       commands,
	}
}

func (app *Application) PublishVideo(ctx context.Context, originalUrl string, savedAddr string) error {
	err := app.PublishDownloadedVideo.Handle(ctx, originalUrl, savedAddr)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) HandleTelegramMessage(ctx context.Context, msg *tgbotapi.Message) error {
	wrappedMessage := message.New(msg)
	for _, ch := range app.TelegramCommands {
		if ch.ShouldRun(wrappedMessage) {
			err := ch.Handle(ctx, wrappedMessage)

			if err != nil {
				return errors.Wrap(err, "telegram command failed")
			}
		}
	}

	return nil
}
