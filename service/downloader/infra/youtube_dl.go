package infra

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/wader/goutubedl"
)

type Downloader struct {
}

func NewDownloader(path string) *Downloader {
	goutubedl.Path = path
	return &Downloader{}
}

func (d *Downloader) DownloadFromUrl(ctx context.Context, url string, w io.Writer) error {
	result, err := goutubedl.New(context.Background(), url, goutubedl.Options{})
	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Infof("start download from %s", url)

	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		logrus.Error(err)
		return err
	}

	defer downloadResult.Close()

	logrus.Infof("start copy %s", url)

	go func() {
		// AK TODO add exception handling
		_, err = io.Copy(w, downloadResult)
	}()

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
