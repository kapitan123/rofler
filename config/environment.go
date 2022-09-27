package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv  = "TELEGRAM_BOT_TOKEN"
	projectIdEnv      = "FIRESTORE_PROJECT_ID"
	gcloudAppCredsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
	port              = "PORT"
	selfUrl           = "SELF_URL"
)

var (
	TelegramToken = os.Getenv(telegramTokenEnv)
	ProjectId     = os.Getenv(projectIdEnv)
	ServerPort, _ = strconv.Atoi(os.Getenv(port))
	GcloudCreds   = os.Getenv(gcloudAppCredsEnv)
	WorkersCount  = 1
	SelfUrl       = os.Getenv(selfUrl)
)

func init() {
	if TelegramToken == "" {
		log.Panic("Telegram bot token is not set. Please set the environment variable ", telegramTokenEnv)
	}

	if ProjectId == "" {
		log.Panic("Firestore projectid is not set. Please set the environment variable ", projectIdEnv)
	}

	if GcloudCreds == "" {
		log.Info("gcloud creds not set. ADC default will be used. Variable ", gcloudAppCredsEnv)
	}

	if SelfUrl == "" {
		log.Info("self url is not set. Bot won't be able to enqueue tasks. Variable ", selfUrl)
	}
}
