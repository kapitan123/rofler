package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TelegramToken          string `env:"TELEGRAM_BOT_TOKEN"`
	VidoFilesBucket        string `env:"VIDEO_FILES_BUCKET"`
	VideoUrlPublishedTopic string `env:"VIDEO_URL_PUBLISHED_TOPIC"`

	ProjectId   string `env:"PROJECT_ID"`
	Port        int    `env:"PORT" envDefault:"8080"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"bot"`

	DebguMode bool `env:"DEBUG_MODE"`
}

func GetEnvVars() Config {
	cfg := Config{}
	err := env.Parse(&cfg)

	if err != nil {
		log.Panic("could not parse env variables")
	}

	if cfg.VidoFilesBucket == "" {
		log.Panic("VIDEO_FILES_BUCKET is not set")
	}

	if cfg.VideoUrlPublishedTopic == "" {
		log.Panic("VIDEO_URL_PUBLISHED_TOPIC is not set")
	}

	if cfg.TelegramToken == "" {
		log.Panic("TELEGRAM_BOT_TOKEN is not set")
	}

	if cfg.ProjectId == "" {
		log.Panic("PROJECT_ID is not set")
	}

	return cfg
}
