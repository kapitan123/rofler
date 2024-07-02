package downloader

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kapitan123/telegrofler/service/downloader/goutubedl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var cookiesName = "cookies.txt"
var ytPath = "yt-dlp"

type Downloader struct {
	cookiesPath string
}

func NewDownloader(base64cookies string) *Downloader {
	cookiesPath, err := createCookieFileFromBase64(base64cookies)

	if err != nil {
		logrus.Panic(errors.Wrap(err, "can't create cookies file"))
	}

	goutubedl.Path = ytPath

	_, err = goutubedl.Version(context.Background())

	if err != nil {
		logrus.Panic(errors.Wrap(err, "can't get yt-dlp version"))
	}

	return &Downloader{cookiesPath}
}

func (d *Downloader) DownloadFromUrl(ctx context.Context, url string, w io.Writer) error {
	result, err := goutubedl.New(context.Background(), url,
		goutubedl.Options{
			CookiesPath: d.cookiesPath,
			DebugLog:    logrus.StandardLogger(),
		})

	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Infof("start download from %s. info: %+v", url, result.Info)

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
		return "", errors.Wrap(err, "error creating local copy of the cookie file")
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", errors.Wrap(err, "error writing to a local copy of the cookie file")
	}

	logrus.Info("cookie file created and content written successfully")

	return getAbsoluteCookiePath()
}

func getAbsoluteCookiePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "couldnt open the current directory")
	}

	cookiesPath := filepath.Join(wd, cookiesName)
	return cookiesPath, nil
}
