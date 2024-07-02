package app

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	VideoFilesBucket  string `env:"VIDEO_FILES_BUCKET"`
	VideoSavedTopicId string `env:"VIDEO_SAVED_TOPIC"`
	DownloaderCookies string `env:"DOWNLOADER_COOKIES"`

	ProjectId   string `env:"PROJECT_ID"`
	Port        int    `env:"PORT" envDefault:"8080"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"downloader"`
	DebguMode   bool   `env:"DEBUG_MODE"`
}

// AK TODO remove panic maybe
func GetEnvVars() Config {
	cfg := Config{}

	err := env.Parse(&cfg)

	if err != nil {
		log.Panic("could not parse env variables")
	}

	if cfg.VideoFilesBucket == "" {
		log.Panic("VIDEO_FILES_BUCKET is not set")
	}

	if cfg.VideoSavedTopicId == "" {
		log.Panic("VIDEO_CONVERTED_TOPIC is not set")
	}

	if cfg.ProjectId == "" {
		log.Panic("PROJECT_ID is not set")
	}

	if cfg.DownloaderCookies == "" {
		log.Panic("DOWNLOADER_COOKIES is not set")
	}

	return cfg
}
