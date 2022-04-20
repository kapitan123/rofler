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
		TelegramToken = "5144696466:AAGPHvo9LGjyqZAUG265hFrpOfn0yjccJ6w"
		//log.Error("Telegram bot token is not set. Please set it in the environment variable ", telegramTokenEnv)
	}

	if ProjectId == "" {
		ProjectId = "neat-cycling-346311"
		// AK TODO panic on all env variables with no defaults
		//log.Panic("Firestore projectid is not set. Please set it in the environment variable ", projectIdEnv)
	}

	if GcloudCreds == "" {
		log.Info("gcloud creds not set. ADC default will be used. Variable ", gcloudAppCredsEnv)
		GcloudCreds = "C:\\pet_proj\\neat-cycling-346311-82a839e49d3d.json"
		// AK TODO panic on all env variables with no defaults
		//log.Panic("Firestore projectid is not set. Please set it in the environment variable ", projectIdEnv)
	}
}
