package infra

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

func (d *Downloader) DownloadFromUrl(url string, w io.Writer) error {
	// youtube-dl -o - "https://www.youtube.com/watch?v=BaW_jenozKcj"
	cmd := exec.Command("youtube-dl", "-o -", url)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	// AK TODO check error handling and also chekc do I even need to spawn it in a separate goroutine.
	// Download may be so slow in the past exactly because container is deprovisioned after serving the request!
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
