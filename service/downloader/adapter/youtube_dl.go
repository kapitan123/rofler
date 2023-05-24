package youtubeDl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Downloader struct {
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (q *Downloader) DownloadFromUrl(url string, w io.Writer) error {
	// youtube-dl -o - "https://www.youtube.com/watch?v=BaW_jenozKcj"
	cmd := exec.Command("youtube-dl", "-o -", url)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	// AK TODO check error handling
	go func() {
		_, err := io.Copy(io.MultiWriter(w, os.Stdout), stdout)
		if err != nil {
			fmt.Println("Error copying output:", err)
		}
	}()

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
