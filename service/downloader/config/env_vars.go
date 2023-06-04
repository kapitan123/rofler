package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	VidoFilesBucket   string `env:"VIDEO_FILES_BUCKET"`
	YoutubeDlPath     string `env:"YOUTUBE_DL_PATH" envDefault:"yt-dlp"`
	VideoSavedTopic   string `env:"VIDEO_SAVED_TOPIC"`
	DownloaderCookies string `env:"DOWNLOADER_COOKIES"`

	ProjectId string `env:"PROJECT_ID"`

	Port        int    `env:"PORT" envDefault:"8080"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"downloader"`
	DebguMode   bool   `env:"DEBUG_MODE"`
}

func GetEnvVars() Config {
	cfg := Config{}

	err := env.Parse(&cfg)

	if err != nil {
		log.Panic("could not parse env variables")
	}

	if cfg.VidoFilesBucket == "" {
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
