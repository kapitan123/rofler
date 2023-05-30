package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TelegramToken      string `env:"TELEGRAM_BOT_TOKEN"`
	VidoFilesBucketUrl string `env:"VIDEO_FILES_BUCKET_URL"`
	DeletionQueueId    string `env:"MESSAGE_DELETION_QUEUE_ID"`
	VideoSavedTopicId  string `env:"VIDEO_SAVED_TOPIC_ID"`
	SelfUrl            string `env:"SELF_URL"`

	ProjectId   string `env:"PROJECT_ID"`
	Port        int    `env:"PORT" envDefault:"8080"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"bot"`

	DebguMode bool `env:"DEBUG_MODE"`
}

func GetEnvVars() Config {
	cfg := Config{}
	env.Parse(&cfg)

	if cfg.VidoFilesBucketUrl == "" {
		log.Panic("VIDEO_FILES_BUCKET_URL is not set")
	}

	if cfg.TelegramToken == "" {
		log.Panic("TELEGRAM_BOT_TOKEN is not set")
	}

	if cfg.SelfUrl == "" {
		log.Panic("SELF_URL is not set")
	}

	if cfg.VideoSavedTopicId == "" {
		log.Panic("VIDEO_SAVED_TOPIC_ID is not set")
	}

	if cfg.DeletionQueueId == "" {
		log.Panic("MESSAGE_DELETION_QUEUE_ID is not set")
	}

	if cfg.ProjectId == "" {
		log.Panic("PROJECT_ID is not set")
	}

	return cfg
}
