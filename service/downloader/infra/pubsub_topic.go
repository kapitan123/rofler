package infra

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
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

func (t *PubSubTopic) PublishSuccess(ctx context.Context, savedVideoAddr string, originalUrl string) error {
	message, _ := json.Marshal(VideoSavedMessage{
		SavedVideoAddr: savedVideoAddr,
		OriginalUrl:    originalUrl,
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

	logrus.Infof("success message for video save was published message id = %s", id)

	return nil
}

type VideoSavedMessage struct {
	SavedVideoAddr string `json:"saved_video_addr"`
	OriginalUrl    string `json:"original_url"`
}
