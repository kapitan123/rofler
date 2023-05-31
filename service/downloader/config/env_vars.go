package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

// AK TODO check if I can make each string a value object type
type Config struct {
	VidoFilesBucketUrl string `env:"VIDEO_FILES_BUCKET_URL" envDefault:"8080"` // AK TODO test
	Port               int    `env:"PORT" envDefault:"8080"`
	VideoSavedTopic    string `env:"VIDEO_SAVED_TOPIC" envDefault:"8080"` // AK TODO test
	ProjectId          string `env:"PROJECT_ID" envDefault:"8080"`
	ServiceName        string `env:"SERVICE_NAME" envDefault:"downloader"`
	DebguMode          bool   `env:"DEBUG_MODE"`
}

func GetEnvVars() Config {
	cfg := Config{}

	err := env.Parse(&cfg)

	if err != nil {
		log.Panic("could not parse env variables")
	}

	if cfg.VidoFilesBucketUrl == "" {
		log.Panic("VIDEO_FILES_BUCKET_URL is not set")
	}

	if cfg.VideoSavedTopic == "" {
		log.Panic("VIDEO_CONVERTED_TOPIC is not set")
	}

	if cfg.ProjectId == "" {
		log.Panic("PROJECT_ID is not set")
	}

	return cfg
}
