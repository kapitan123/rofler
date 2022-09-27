// Should not be used outside of main.go
// will be converted to flags
package config

import (
	"os"
	"strconv"

	"cloud.google.com/go/compute/metadata"
	log "github.com/sirupsen/logrus"
)

const (
	telegramTokenEnv  = "TELEGRAM_BOT_TOKEN"
	gcloudAppCredsEnv = "GOOGLE_APPLICATION_CREDENTIALS"
	port              = "PORT"
	selfUrl           = "SELF_URL"
)

var (
	TelegramToken = os.Getenv(telegramTokenEnv)
	ServerPort, _ = strconv.Atoi(os.Getenv(port))
	GcloudCreds   = os.Getenv(gcloudAppCredsEnv)
	WorkersCount  = 1
	SelfUrl       = os.Getenv(selfUrl)
)

type Meta struct {
	ProjectId string
	Region    string
	Email     string
}

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

func GetMetadata() (*Meta, error) {
	projectId, err := metadata.ProjectID()

	if err != nil {
		return nil, err
	}

	region, err := metadata.Zone()

	if err != nil {
		return nil, err
	}

	email, err := metadata.Get("instance/service-accounts/default/email")

	if err != nil {
		return nil, err
	}

	return &Meta{
		ProjectId: projectId,
		Region:    region,
		Email:     email,
	}, nil
}
