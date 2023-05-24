// Should not be used outside of main.go
// will be converted to flags
package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv           = "TELEGRAM_BOT_TOKEN"
	gcloudAppCredsEnv          = "GOOGLE_APPLICATION_CREDENTIALS" // Default cloudrun env -should be set to run local container
	port                       = "PORT"
	selfUrl                    = "SELF_URL"
	region                     = "REGION"
	saEmail                    = "SA_EMAIL"
	projectId                  = "PROJECT_ID"
	serviceName                = "K_SERVICE" // Default cloudrun env
	videoWasPublishedTopicName = "VIDEO_PUBLISHED_TOPIC"
	deletionQueueId            = "MESSAGE_DELETION_QUEUE_ID"
)

var (
	TelegramToken              = os.Getenv(telegramTokenEnv)
	ServerPort, _              = strconv.Atoi(os.Getenv(port))
	WorkersCount               = 1
	DeletionQueueId            = os.Getenv(deletionQueueId)
	VideoWasPublishedTopicName = os.Getenv(videoWasPublishedTopicName)
)

// these var are used to manually set context, which is fetched from gcloud metadata when deployed
var (
	SelfUrl     = os.Getenv(selfUrl)
	Region      = os.Getenv(region)
	SaEmail     = os.Getenv(saEmail)
	ServiceName = os.Getenv(serviceName)
	ProjectId   = os.Getenv(projectId)
)

func init() {
	if TelegramToken == "" {
		log.Panic("telegram bot token is not set. Please set the environment variable ", telegramTokenEnv)
	}

	if SelfUrl == "" {
		log.Info("self url is not set. Bot won't be able to enqueue tasks. Variable ", selfUrl)
	}

	if SaEmail == "" || Region == "" || ProjectId == "" || SelfUrl == "" {
		log.Infof("metadata variables are not set, metadata server will be used %s, %s, %s ", saEmail, region, projectId)
	}
}
