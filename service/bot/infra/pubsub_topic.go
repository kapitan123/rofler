package infra

import (
	"context"
	"encoding/json"
	"net/url"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

type PubSubTopic struct {
	topic  *pubsub.Topic
	origin string
}

func NewPubSubTopicClient(ctx context.Context, projectId string, servicename string, videoUrlPublishedTopicId string) *PubSubTopic {
	newPubSubClient, err := pubsub.NewClient(ctx, projectId)

	if err != nil {
		panic(err)
	}

	return &PubSubTopic{
		topic:  newPubSubClient.Topic(videoUrlPublishedTopicId),
		origin: servicename,
	}
}

func (t *PubSubTopic) PublishUrl(ctx context.Context, url *url.URL) error {
	message, _ := json.Marshal(VideoUrlPostedMessage{
		Url: url.String(),
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

	logrus.Infof("videoUrlPostedMessage was published, message id = %s", id)

	return nil
}

type VideoUrlPostedMessage struct {
	Url string `json:"url"`
}
