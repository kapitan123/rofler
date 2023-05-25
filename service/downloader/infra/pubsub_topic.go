package adapter

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PubSubTopic struct {
	topic  *pubsub.Topic
	origin string
}

func NewPubSubTopicClient(ctx context.Context, projectId string, servicename string, videoSavedTopicId string) *PubSubTopic {
	newPubSubClient, err := pubsub.NewClient(ctx, projectId)

	if err != nil {
		panic(err)
	}

	return &PubSubTopic{
		topic:  newPubSubClient.Topic(videoSavedTopicId),
		origin: servicename,
	}
}

func (t *PubSubTopic) PublishSuccess(ctx context.Context, savedVideoId uuid.UUID, originalUrl string) error {
	message, _ := json.Marshal(VideoSavedMessage{
		SavedVideoId: savedVideoId,
		OriginalUrl:  originalUrl,
	})
	result := t.topic.Publish(ctx, &pubsub.Message{
		Data: message,
		Attributes: map[string]string{
			"origin": t.origin,
		},
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)

	if err != nil {
		return err
	}

	logrus.Info("success message for video save was published message id = %s", id)

	return nil
}

type VideoSavedMessage struct {
	SavedVideoId uuid.UUID `json:"saved_video_id"`
	OriginalUrl  string    `json:"original_url"`
}
