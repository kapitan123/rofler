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

func NewPubSubTopicClient(ctx context.Context, projectId string, servicename string, videoUrlPostedTopicId string) *PubSubTopic {
	newPubSubClient, err := pubsub.NewClient(ctx, projectId)

	if err != nil {
		panic(err)
	}

	return &PubSubTopic{
		topic:  newPubSubClient.Topic(videoUrlPostedTopicId),
		origin: servicename,
	}
}

func (t *PubSubTopic) PublishUrl(ctx context.Context, url string) error {
	message, _ := json.Marshal(VideoUrlPostedMessage{
		Url: url,
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

	logrus.Info("videoUrlPostedMessage was published, message id = %s", id)

	return nil
}

type VideoUrlPostedMessage struct {
	Url string `json:"url"`
}
