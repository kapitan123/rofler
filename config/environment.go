// Should not be used outside of main.go
// will be converted to flags
package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv  = "TELEGRAM_BOT_TOKEN"
	gcloudAppCredsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
	port              = "PORT"
	selfUrl           = "SELF_URL"
	deletionQueueName = "DELETION_QUEUE_NAME"
)

var (
	TelegramToken     = os.Getenv(telegramTokenEnv)
	ServerPort, _     = strconv.Atoi(os.Getenv(port))
	GcloudCreds       = os.Getenv(gcloudAppCredsEnv)
	WorkersCount      = 1
	SelfUrl           = os.Getenv(selfUrl)
	DeletionQueueName = os.Getenv(deletionQueueName)
)

func init() {
	if TelegramToken == "" {
		log.Panic("telegram bot token is not set. Please set the environment variable ", telegramTokenEnv)
	}

	if GcloudCreds == "" {
		log.Info("gcloud creds not set. ADC default will be used. Variable ", gcloudAppCredsEnv)
	}

	if SelfUrl == "" {
		log.Info("self url is not set. Bot won't be able to enqueue tasks. Variable ", selfUrl)
	}
}
