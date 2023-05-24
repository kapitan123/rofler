package app

import (
	"io"
)

type Application struct {
	videoSavedTopic  topic
	videoFilesBucket fileBucket
	downloader       downloader
}

type topic interface {
	PublishSuccess(savedVideoAddr string, originalUrl string) error
}

type fileBucket interface {
	Save(fromReader io.Reader) (string, error)
	Read(addr string, r io.Reader) error
}

type downloader interface {
	DownloadFromUrl(url string, w io.Writer) error
}

func NewApplication(videoSavedTopic topic, videoBucket fileBucket, downloader downloader) Application {
	return Application{
		videoSavedTopic:  videoSavedTopic,
		videoFilesBucket: videoBucket,
		downloader:       downloader,
	}
}

func (app *Application) SaveVideoToStorage(url string) error {
	pipeReader, pipeWriter := io.Pipe()

	err := app.downloader.DownloadFromUrl(url, pipeWriter)

	if err != nil {
		return err
	}

	path, err := app.videoFilesBucket.Save(pipeReader)

	if err != nil {
		return err
	}

	err = app.videoSavedTopic.PublishSuccess(path, url)

	if err != nil {
		return err
	}

	return err
}

func (app *Application) GetVideo() {
	// AK TODO implement
}
