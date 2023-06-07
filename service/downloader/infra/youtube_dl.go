package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wader/goutubedl"
)

var cookiesName = "cookies.txt"

type Downloader struct {
}

func NewDownloader(path string, base64cookies string) *Downloader {
	cookiesPath, err := createCookieFileFromBase64(base64cookies)

	if err != nil {
		logrus.Panic(err)
	}

	goutubedl.Path = fmt.Sprintf("%s --cookies %s", path, cookiesPath)

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

	defer func() {
		logrus.Infof("finish download from %s", url)
		downloadResult.Close()
	}()

	_, err = io.Copy(w, downloadResult)

	if err != nil {
		return err
	}

	return nil
}

func createCookieFileFromBase64(base64cookies string) (string, error) {
	decodedDataBytes, err := base64.StdEncoding.DecodeString(base64cookies)
	content := string(decodedDataBytes)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("couldn't parse cookies websites requiring authentication will not be accessible %s", base64cookies))
	}

	file, err := os.Create(cookiesName)
	if err != nil {
		return "", errors.Wrap(err, "error creating file")
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", errors.Wrap(err, "error writing to file")
	}

	logrus.Info("cookie file created and content written successfully")

	return getAbsoluteCookiePath()
}

func getAbsoluteCookiePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cookiesPath := filepath.Join(wd, cookiesName)
	return cookiesPath, nil
}
