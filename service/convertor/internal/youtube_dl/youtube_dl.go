package youtubeDl

import (
	"fmt"
	"io"
	"os/exec"

	"cloud.google.com/go/storage"
)

type YoutubeDl struct {
	client *storage.Client
}

func New() *YoutubeDl {
	return &YoutubeDl{}
}

func (q *YoutubeDl) DownloadUrl(url string, w io.Writer) error {
	// youtube-dl -o - "https://www.youtube.com/watch?v=BaW_jenozKcj"
	cmd := exec.Command("youtube-dl", "-o -", url)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
