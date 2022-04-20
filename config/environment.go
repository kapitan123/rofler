package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv  = "TELEGRAM_BOT_TOKEN"
	projectIdEnv      = "FIRESTORE_PROJECT_ID"
	gcloudAppCredsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
)

var TelegramToken = os.Getenv(telegramTokenEnv)
var ProjectId = os.Getenv(projectIdEnv)
var ServerPort = 9001
var GcloudCreds = os.Getenv(gcloudAppCredsEnv)

func init() {
	if TelegramToken == "" {
		log.Panic("Telegram bot token is not set. Please set it in the environment variable ", telegramTokenEnv)
	}

	if ProjectId == "" {
		ProjectId = "neat-cycling-346311"
		log.Panic("Firestore projectid is not set. Please set it in the environment variable ", projectIdEnv)
	}

	if GcloudCreds == "" {
		log.Info("gcloud creds not set. ADC default will be used. Variable ", gcloudAppCredsEnv)
	}
}
